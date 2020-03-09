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
	"context"
	"encoding/hex"
	fmt "fmt"
	"log"
	"os"
	"time"

	"github.com/lnsp/ftp2p/pkg/hash"
	"github.com/lnsp/ftp2p/pkg/security"
	"github.com/lnsp/ftp2p/pkg/tavern"
	"github.com/lnsp/ftp2p/pkg/tracker"
	"github.com/lnsp/grpc-quic/opts"

	qgrpc "github.com/lnsp/grpc-quic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	conn, err := qgrpc.Dial(tracker.Addr, opts.WithInsecure(), opts.WithBlock(), opts.WithTimeout(time.Minute), opts.WithBackoffMaxDelay(time.Minute))
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
	// Verify hash for tracker
	if err := hash.Verify(path, tracker.Hash); err != nil {
		return fmt.Errorf("verify path: %w", err)
	}
	log.Printf("Add tracker %s for %s", hex.EncodeToString(tracker.Hash), path)
	seeder.index.Add(path, tracker)
	go func() {
		// Announce to tavern
		for {
			interval, err := seeder.Announce(tracker)
			if err != nil {
				log.Printf("Announce to tavern: %v", err)
				time.Sleep(time.Minute)
			} else {
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}
	}()
	return nil
}

func (seeder *Seeder) ListenAndServe() error {
	tlsCfg, err := security.NewBasicTLS()
	if err != nil {
		return err
	}
	grpcServer, listener, err := qgrpc.NewServer(seeder.Addr, opts.TLSConfig(tlsCfg))
	if err != nil {
		return err
	}
	RegisterSeederServer(grpcServer, seeder)
	return grpcServer.Serve(listener)
}

func New(addr string, poolsize int) *Seeder {
	return &Seeder{
		Addr:   addr,
		fdpool: newPool(poolsize),
		index:  newFileIndex(),
	}
}
