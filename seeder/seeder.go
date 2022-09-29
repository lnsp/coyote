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
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	seederv1 "github.com/lnsp/ftp2p/gen/seeder/v1"
	"github.com/lnsp/ftp2p/gen/seeder/v1/seederv1connect"
	tavernv1 "github.com/lnsp/ftp2p/gen/tavern/v1"
	"github.com/lnsp/ftp2p/gen/tavern/v1/tavernv1connect"
	trackerv1 "github.com/lnsp/ftp2p/gen/tracker/v1"
	"github.com/lnsp/ftp2p/hash"
	"github.com/lnsp/ftp2p/http3utils"
	"github.com/lucas-clemente/quic-go/http3"

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
	Addr     string
	Insecure bool

	index     *FileIndex
	fdpool    *Pool
	tlsConfig *tls.Config
}

func (seeder *Seeder) Has(ctx context.Context, req *connect.Request[seederv1.HasRequest]) (*connect.Response[seederv1.HasResponse], error) {
	chunks := seeder.index.Scan(req.Msg.Hash)
	return connect.NewResponse(&seederv1.HasResponse{
		Chunks: chunks.Bytes(),
	}), nil
}

func (seeder *Seeder) Fetch(ctx context.Context, req *connect.Request[seederv1.FetchRequest]) (*connect.Response[seederv1.FetchResponse], error) {
	chunk, ok := seeder.index.Get(req.Msg.Hash, req.Msg.Chunk)
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
	return connect.NewResponse(&seederv1.FetchResponse{
		Data: buf,
	}), nil
}

const minAnnounceBackoff = time.Minute

func (seeder *Seeder) Announce(ctx context.Context, tracker *trackerv1.Tracker, addr string) error {
	httpClient := http3utils.DefaultClient
	if seeder.Insecure {
		httpClient = http3utils.DefaultClientInsecure
	}
	client := tavernv1connect.NewTavernServiceClient(httpClient, addr)
	// First, we need to parse our local certificate
	if len(seeder.tlsConfig.Certificates) > 1 || len(seeder.tlsConfig.Certificates[0].Certificate) > 1 {
		return fmt.Errorf("can only handle single seeder certificate")
	}
	certificate, err := x509.ParseCertificate(seeder.tlsConfig.Certificates[0].Certificate[0])
	if err != nil {
		return fmt.Errorf("failed to parse seeder certificate")
	}
	publicKey, err := x509.MarshalPKIXPublicKey(certificate.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal seeder public key")
	}

	sleep := minAnnounceBackoff
	timer := time.NewTimer(0)
	for {
		// Submit hash, seeder address and certificate to tavern
		resp, err := client.Announce(ctx, connect.NewRequest(&tavernv1.AnnounceRequest{
			Hash:      tracker.Hash,
			Addr:      "https://" + seeder.Addr,
			PublicKey: publicKey,
		}))
		if err != nil {
			log.Printf("Announce to Tavern %s: %v", addr, err)
			sleep = minAnnounceBackoff
		} else {
			sleep = time.Duration(resp.Msg.Interval) * time.Second
		}
		// Reset announcement timer
		timer.Reset(sleep)
		select {
		case <-timer.C:
		case <-ctx.Done():
			return nil
		}
	}
}

// Seed verifies that the given tracker matches the file at the path and begins announcing seeding capacities
// to the taverns listed in the tracker. All background routines keep running until the context is canceled.
func (seeder *Seeder) Seed(ctx context.Context, path string, tracker *trackerv1.Tracker) error {
	// Verify hash for tracker
	if err := hash.Verify(path, tracker.Hash); err != nil {
		return fmt.Errorf("verify path: %w", err)
	}
	log.Printf("Add tracker %s for %s", hex.EncodeToString(tracker.Hash), path)
	seeder.index.Add(path, tracker)
	// Start announcing to each Tavern mentioned in the Tracker
	for _, addr := range tracker.Taverns {
		// Announce to tavern
		go seeder.Announce(ctx, tracker, addr)
	}
	return nil
}

func (seeder *Seeder) ListenAndServe(ctx context.Context) error {
	// Set up handler and mux
	path, handler := seederv1connect.NewSeederServiceHandler(seeder)
	mux := http.NewServeMux()
	mux.Handle(path, handler)

	// Setup HTTP3 server
	server := http3.Server{
		Addr:      seeder.Addr,
		TLSConfig: seeder.tlsConfig,
		Handler:   handler,
	}
	go func() { <-ctx.Done(); server.Close() }()
	return server.ListenAndServe()
}

func New(addr string, insecure bool, poolsize int, tlsConfig *tls.Config) *Seeder {
	// Setup TLS configuration
	return &Seeder{
		Addr:      addr,
		Insecure:  insecure,
		fdpool:    newPool(poolsize),
		index:     newFileIndex(),
		tlsConfig: tlsConfig,
	}
}
