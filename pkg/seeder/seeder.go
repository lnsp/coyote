// Copyright 2020 Lennart Espe
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package seeder

import (
	"bytes"
	"context"
	"encoding/hex"
	fmt "fmt"
	"github.com/lnsp/ftp2p/pkg/tavern"
	"github.com/lnsp/ftp2p/pkg/tracker"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/willf/bitset"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Pool struct {
	buffer chan bool
}

func (pool *Pool) Get(ctx context.Context) bool {
	select {
	case <-pool.buffer:
		return true
	case <-ctx.Done():
		return false
	case <-time.After(time.Minute):
		return false
	}
}

func (pool *Pool) Put() {
	pool.buffer <- true
}

func newPool(size int) *Pool {
	p := make(chan bool, size)
	for i := 0; i < size; i++ {
		p <- true
	}
	return &Pool{p}
}

func byteSliceComparator(a, b interface{}) int {
	return bytes.Compare(a.([]byte), b.([]byte))
}

type FileEntry struct {
	Path   string
	Chunk  int64
	Offset int64
	Size   int64
}

type FileIndex struct {
	sync.RWMutex
	treemap *treemap.Map
}

func (index *FileIndex) Scan(hash []byte) *bitset.BitSet {
	index.RLock()
	value, ok := index.treemap.Get(hash)
	index.RUnlock()
	if !ok {
		return bitset.New(0)
	}
	chunks := bitset.New(uint(tracker.MinChunkCount))
	entries := value.([]*FileEntry)
	for _, e := range entries {
		chunks.Set(uint(e.Chunk))
	}
	return chunks
}

func (index *FileIndex) Get(hash []byte, chunk int64) (*FileEntry, bool) {
	index.RLock()
	value, ok := index.treemap.Get(hash)
	index.RUnlock()
	if !ok {
		return nil, false
	}
	entries := value.([]*FileEntry)
	found := sort.Search(len(entries), func(i int) bool {
		return entries[i].Chunk >= chunk
	})
	if found >= len(entries) || entries[found].Chunk != chunk {
		return nil, false
	}
	return entries[found], true
}

func (index *FileIndex) Add(path string, tracker *tracker.Tracker) {
	var offset int64
	entries := make([]*FileEntry, len(tracker.ChunkHashes))
	for chunk := range tracker.ChunkHashes {
		size := tracker.ChunkSize
		if chunk == len(tracker.ChunkHashes)-1 {
			size = tracker.Size - offset
		}
		entries[chunk] = &FileEntry{
			Path:   path,
			Chunk:  int64(chunk),
			Offset: offset,
			Size:   size,
		}
		offset += size
	}
	index.Lock()
	index.treemap.Put(tracker.Hash, entries)
	index.Unlock()
}

func newFileIndex() *FileIndex {
	return &FileIndex{
		treemap: treemap.NewWith(byteSliceComparator),
	}
}

type Seeder struct {
	Addr   string
	index  *FileIndex
	fdpool *Pool
}

func (seeder *Seeder) Has(ctx context.Context, req *HasRequest) (*HasResponse, error) {
	chunks := seeder.index.Scan(req.Hash)
	return &HasResponse{
		Chunks: chunks.Bytes(),
	}, nil
}

func (seeder *Seeder) Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error) {
	chunk, ok := seeder.index.Get(req.Hash, req.Chunk)
	if !ok {
		return nil, status.Error(codes.NotFound, "chunk not found")
	}
	if !seeder.fdpool.Get(ctx) {
		return nil, status.Error(codes.Canceled, "request canceled")
	}
	defer seeder.fdpool.Put()
	chunkfile, err := os.Open(chunk.Path)
	if err != nil {
		return nil, status.Error(codes.Internal, "file access failure")
	}
	defer chunkfile.Close()
	buf := make([]byte, chunk.Size)
	if _, err := chunkfile.ReadAt(buf, chunk.Offset); err != nil {
		return nil, status.Error(codes.Internal, "file read failure")
	}
	return &FetchResponse{
		Data: buf,
	}, nil
}

func (seeder *Seeder) Announce(tracker *tracker.Tracker) (int64, error) {
	conn, err := grpc.Dial(tracker.Addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Minute), grpc.WithBackoffMaxDelay(time.Minute))
	if err != nil {
		return 0, fmt.Errorf("announce to tavern: %v", err)
	}
	defer conn.Close()
	client := tavern.NewTavernClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	resp, err := client.Announce(ctx, &tavern.AnnounceRequest{
		Hash: tracker.Hash,
		Addr: seeder.Addr,
	})
	if err != nil {
		return 0, fmt.Errorf("announce to tavern: %v", err)
	}
	return resp.Interval, nil
}

func (seeder *Seeder) Seed(path string, tracker *tracker.Tracker) error {
	log.Printf("Add tracker %s for %s", hex.EncodeToString(tracker.Hash), path)
	seeder.index.Add(path, tracker)
	// Announce to tavern
	_, err := seeder.Announce(tracker)
	if err != nil {
		return fmt.Errorf("seeder announce tracker: %v", err)
	}
	return nil
}

func (seeder *Seeder) ListenAndServe() error {
	listener, err := net.Listen("tcp", seeder.Addr)
	if err != nil {
		return fmt.Errorf("seeder listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterSeederServer(grpcServer, seeder)
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("seeder serve: %v", err)
	}
	return nil
}

func New(addr string, poolsize int) *Seeder {
	return &Seeder{
		Addr:   addr,
		fdpool: newPool(poolsize),
		index:  newFileIndex(),
	}
}
