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
	context "context"
	"crypto/tls"
	"encoding/hex"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/lucas-clemente/quic-go/http3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tavernv1 "github.com/lnsp/ftp2p/gen/tavern/v1"
	"github.com/lnsp/ftp2p/gen/tavern/v1/tavernv1connect"
)

type Tavern struct {
	AnnounceInterval int64
	index            *PeerIndex

	cleanupLock  sync.Mutex
	cleanupQueue [][]byte

	tavernv1connect.UnimplementedTavernServiceHandler
}

// List all announced peers for the given hash.
func (tavern *Tavern) List(ctx context.Context, req *connect.Request[tavernv1.ListRequest]) (*connect.Response[tavernv1.ListResponse], error) {
	log.Printf("List peers for %s", hex.EncodeToString(req.Msg.Hash))
	entries, ok := tavern.index.Find(req.Msg.Hash)
	if !ok {
		return nil, status.Error(codes.NotFound, "hash not found")
	}
	peers := make([]*tavernv1.Peer, len(entries))
	expired := false
	for i, e := range entries {
		if e.Timestamp.Add(time.Duration(tavern.AnnounceInterval) * time.Second).Before(time.Now()) {
			expired = true
		}
		peers[i] = e.Peer
	}
	if expired {
		tavern.scheduleForCleanup(req.Msg.Hash, 0)
	}
	return connect.NewResponse(&tavernv1.ListResponse{
		Peers:    peers,
		Interval: tavern.AnnounceInterval,
	}), nil
}

// Announce a new peer for the given hash.
func (tavern *Tavern) Announce(ctx context.Context, req *connect.Request[tavernv1.AnnounceRequest]) (*connect.Response[tavernv1.AnnounceResponse], error) {
	log.Printf("Announce peer %s for %s", req.Msg.Addr, hex.EncodeToString(req.Msg.Hash))

	tavern.index.Seed(req.Msg.Hash, &tavernv1.Peer{
		Addr:      req.Msg.Addr,
		PublicKey: req.Msg.PublicKey,
	})
	tavern.scheduleForCleanup(req.Msg.Hash, time.Duration(tavern.AnnounceInterval)*time.Second)
	return connect.NewResponse(&tavernv1.AnnounceResponse{
		Interval: tavern.AnnounceInterval,
	}), nil
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

func (tavern *Tavern) ListenAndServe(addr string, tlsConfig *tls.Config) error {
	path, handler := tavernv1connect.NewTavernServiceHandler(tavern)
	// Setup mux
	mux := http.NewServeMux()
	mux.Handle(path, handler)
	// Setup Connect handler
	server := &http3.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
		Handler:   mux,
	}
	return server.ListenAndServe()
}

func New(interval int64) *Tavern {
	return &Tavern{
		AnnounceInterval: interval,
		index:            newPeerIndex(),
	}
}
