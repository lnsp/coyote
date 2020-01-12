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

package tavern

import (
	"bytes"
	context "context"
	"encoding/hex"
	fmt "fmt"
	"log"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PeerEntry struct {
	Peer      *Peer
	Timestamp time.Time
}

type PeerIndex struct {
	sync.RWMutex
	treemap *treemap.Map
}

func newPeerIndex() *PeerIndex {
	return &PeerIndex{
		treemap: treemap.NewWith(byteSliceComparator),
	}
}

func (index *PeerIndex) Find(hash []byte) ([]PeerEntry, bool) {
	index.RLock()
	value, ok := index.treemap.Get(hash)
	index.RUnlock()
	if !ok {
		return nil, false
	}
	return value.([]PeerEntry), true
}

func (index *PeerIndex) Clean(hash []byte, timeout int64) {
	index.Lock()
	defer index.Unlock()
	value, ok := index.treemap.Get(hash)
	if !ok {
		return
	}
	entries := value.([]PeerEntry)
	n := sort.Search(len(entries), func(i int) bool {
		return entries[i].Timestamp.After(time.Now().Add(time.Duration(-timeout) * time.Second))
	})
	index.treemap.Put(hash, entries[n:])
}

func (index *PeerIndex) Seed(hash []byte, peer *Peer) {
	index.Lock()
	value, ok := index.treemap.Get(hash)
	var entries []PeerEntry
	if ok {
		entries = value.([]PeerEntry)
	}
	entries = append(entries, PeerEntry{peer, time.Now()})
	index.treemap.Put(hash, entries)
	index.Unlock()
}

func byteSliceComparator(a, b interface{}) int {
	return bytes.Compare(a.([]byte), b.([]byte))
}

type Tavern struct {
	AnnounceInterval int64
	index            *PeerIndex

	cleanupLock  sync.Mutex
	cleanupQueue [][]byte
}

func (tavern *Tavern) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	log.Printf("List peers for %s", hex.EncodeToString(req.Hash))
	entries, ok := tavern.index.Find(req.Hash)
	if !ok {
		return nil, status.Error(codes.NotFound, "hash not found")
	}
	peers := make([]*Peer, len(entries))
	expired := false
	for i, e := range entries {
		if e.Timestamp.Add(time.Duration(tavern.AnnounceInterval) * time.Second).Before(time.Now()) {
			expired = true
		}
		peers[i] = e.Peer
	}
	if expired {
		tavern.scheduleForCleanup(req.Hash, 0)
	}
	return &ListResponse{
		Peers:    peers,
		Interval: tavern.AnnounceInterval,
	}, nil
}

func (tavern *Tavern) Announce(ctx context.Context, req *AnnounceRequest) (*AnnounceResponse, error) {
	log.Printf("Announce peer %s for %s", req.Addr, hex.EncodeToString(req.Hash))
	tavern.index.Seed(req.Hash, &Peer{
		Addr: req.Addr,
	})
	tavern.scheduleForCleanup(req.Hash, time.Duration(tavern.AnnounceInterval)*time.Second)
	return &AnnounceResponse{
		Interval: tavern.AnnounceInterval,
	}, nil
}

func (tavern *Tavern) scheduleForCleanup(hash []byte, after time.Duration) {
	go func() {
		log.Printf("Scheduled %s for cleanup after %s", hex.EncodeToString(hash), after)
		time.Sleep(after)
		tavern.cleanupLock.Lock()
		tavern.cleanupQueue = append(tavern.cleanupQueue, hash)
		tavern.cleanupLock.Unlock()
	}()
}

func (tavern *Tavern) cleanupWorker() {
	for {
		if len(tavern.cleanupQueue) < 1 {
			time.Sleep(time.Second)
			continue
		}
		var hash []byte
		tavern.cleanupLock.Lock()
		hash, tavern.cleanupQueue = tavern.cleanupQueue[0], tavern.cleanupQueue[1:]
		tavern.cleanupLock.Unlock()
		log.Printf("Cleaning up %s", hex.EncodeToString(hash))
		tavern.index.Clean(hash, tavern.AnnounceInterval)
	}
}

func (tavern *Tavern) ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("tcp listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterTavernServer(grpcServer, tavern)
	log.Printf("Hosting tavern on %s", addr)

	go tavern.cleanupWorker()
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("tavern serve: %v", err)
	}
	return nil
}

func New() *Tavern {
	return &Tavern{
		AnnounceInterval: 300,
		index:            newPeerIndex(),
	}
}
