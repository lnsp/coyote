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
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PeerIndex struct {
	sync.RWMutex
	treemap *treemap.Map
}

func newPeerIndex() *PeerIndex {
	return &PeerIndex{
		treemap: treemap.NewWith(byteSliceComparator),
	}
}

func (index *PeerIndex) Find(hash []byte) ([]*Peer, bool) {
	index.RLock()
	value, ok := index.treemap.Get(hash)
	index.RUnlock()
	if !ok {
		return nil, false
	}
	return value.([]*Peer), true
}

func (index *PeerIndex) Seed(hash []byte, peer *Peer) {
	index.Lock()
	value, ok := index.treemap.Get(hash)
	var peerlist []*Peer
	if ok {
		peerlist = append(value.([]*Peer), peer)
	} else {
		peerlist = []*Peer{peer}
	}
	index.treemap.Put(hash, peerlist)
	index.Unlock()
}

func byteSliceComparator(a, b interface{}) int {
	return bytes.Compare(a.([]byte), b.([]byte))
}

type Tavern struct {
	AnnounceInterval int64
	index            *PeerIndex
}

func (tavern *Tavern) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	log.Printf("List peers for %s", hex.EncodeToString(req.Hash))
	peers, ok := tavern.index.Find(req.Hash)
	if !ok {
		return nil, status.Error(codes.NotFound, "hash not found")
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
	return &AnnounceResponse{
		Interval: tavern.AnnounceInterval,
	}, nil
}

func (tavern *Tavern) ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("tcp listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterTavernServer(grpcServer, tavern)
	log.Printf("Hosting tavern on %s", addr)
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
