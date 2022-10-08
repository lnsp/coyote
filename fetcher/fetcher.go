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

package fetcher

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	seederv1 "github.com/lnsp/coyote/gen/seeder/v1"
	"github.com/lnsp/coyote/gen/seeder/v1/seederv1connect"
	tavernv1 "github.com/lnsp/coyote/gen/tavern/v1"
	"github.com/lnsp/coyote/gen/tavern/v1/tavernv1connect"
	trackerv1 "github.com/lnsp/coyote/gen/tracker/v1"
	"github.com/lnsp/coyote/http3utils"

	"github.com/bits-and-blooms/bitset"
	"github.com/schollz/progressbar/v3"
)

// ListPeers contacts a tracker looking for a given hash.
func (fetcher *Fetcher) ListPeers(hash []byte, addr string) ([]Peer, error) {
	tavernClient := tavernv1connect.NewTavernServiceClient(fetcher.httpClient, addr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	resp, err := tavernClient.List(ctx, connect.NewRequest(&tavernv1.ListRequest{
		Hash: hash,
	}))
	if err != nil {
		return nil, fmt.Errorf("list peers: %v", err)
	}
	peers := make([]Peer, len(resp.Msg.Peers))
	for i := range resp.Msg.Peers {
		// Generate HTTP client for Peer
		httpClient, err := http3utils.TrustedClient(resp.Msg.Peers[i].PublicKey)
		if err != nil {
			return nil, fmt.Errorf("create client for peer: %w", err)
		}
		peerClient := seederv1connect.NewSeederServiceClient(httpClient, resp.Msg.Peers[i].Addr)
		peers[i] = Peer{
			Address: resp.Msg.Peers[i].Addr,
			Hash:    hash,
			Client:  peerClient,
		}
	}
	return peers, nil
}

// HasChunks fetches the set of chunks served by the peer.
func (peer *Peer) HasChunks(hash []byte) (*bitset.BitSet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	resp, err := peer.Client.Has(ctx, connect.NewRequest(&seederv1.HasRequest{
		Hash: hash,
	}))
	if err != nil {
		return nil, fmt.Errorf("seeder has: %v", err)
	}
	return bitset.From(resp.Msg.Chunks), nil
}

type Peer struct {
	Address  string
	Hash     []byte
	Chunkset *bitset.BitSet
	Client   seederv1connect.SeederServiceClient
}

func (peer *Peer) Refresh() {
	chunkset, err := peer.HasChunks(peer.Hash)
	if err != nil {
		log.Printf("fetch chunkset from peer %s: %v", peer.Address, err)
		peer.Chunkset = bitset.New(0)
	} else {
		peer.Chunkset = chunkset
	}
}

type Task struct {
	Fetched     chan int64
	Tracker     *trackerv1.Tracker
	Peers       []Peer
	Peer        int
	Hash        []byte
	ChunkHash   []byte
	Chunk       int64
	Destination string
}

func (task Task) Fetch() (int64, error) {
	// Use client from pool
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	resp, err := task.Peers[task.Peer].Client.Fetch(ctx, connect.NewRequest(&seederv1.FetchRequest{
		Hash:  task.Hash,
		Chunk: task.Chunk,
	}))
	if err != nil {
		return -1, fmt.Errorf("seeder fetch: %v", err)
	}
	hasher := sha256.New()
	hasher.Write(resp.Msg.Data)
	if !bytes.Equal(hasher.Sum(nil), task.ChunkHash) {
		return -1, fmt.Errorf("bad chunk from peer")
	}
	if err := os.WriteFile(task.Destination, resp.Msg.Data, 0644); err != nil {
		return -1, fmt.Errorf("chunk write: %v", err)
	}
	return int64(len(resp.Msg.Data)), nil
}

func generateChunkPath(basepath string, i int) string {
	return fmt.Sprintf("%s.%d", basepath, i)
}

func randomIntOrder(n int) []int {
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	rand.Shuffle(n, func(i, j int) { order[i], order[j] = order[j], order[i] })
	return order
}

func (fetcher *Fetcher) handleFetchTasks(tasks chan Task) {
	for task := range tasks {
		chunk := uint(task.Chunk)
		fetched := int64(-1)
		// Check for local correct copy
		if data, err := readAndVerify(task.Destination, task.ChunkHash); err == nil {
			fetched = int64(len(data))
		}
		// If not exist, try to fetch
		for cycle := 0; fetched < 0; cycle++ {
			// To load balance, generate random peer order per cycle
			order := randomIntOrder(len(task.Peers))
			for i := 0; i < len(task.Peers) && fetched < 0; i++ {
				j := order[i]
				if !task.Peers[j].Chunkset.Test(chunk) {
					log.Printf("peer %s does not have chunk %d", task.Peers[j].Address, chunk)
					continue
				}
				task.Peer = j
				n, err := task.Fetch()
				if err != nil {
					log.Printf("fetch chunk %d from peer %s failed: %v", chunk, task.Peers[j].Address, err)
					continue
				}
				fetched = n
			}
			// If fetched is still -1, we need to re-fetch the peer list
			if fetched < 0 {
				log.Printf("worker cycled for chunk %d in iteration %d", chunk, cycle)
				// Do exponential backoff, reaches max around 16 cycles
				backoff := time.Second * time.Duration(math.Min(60.0, math.Pow(1.3, float64(cycle))))
				time.Sleep(backoff)
				// To resolve worker cycle, re-fetch peers and chunksets
				peers, err := fetcher.Resolve(task.Tracker)
				if err != nil {
					log.Printf("resolve after cycle: %v", err)
				} else {
					task.Peers = peers
				}
			}
		}
		// Submit success
		task.Fetched <- fetched
	}
}

func readAndVerify(path string, hash []byte) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read and verify: %w", err)
	}
	hasher := sha256.New()
	hasher.Write(data)
	sum := hasher.Sum(nil)
	if !bytes.Equal(sum, hash) {
		return nil, fmt.Errorf("hash does not match")
	}
	return data, nil
}

type Fetcher struct {
	maxConns   int
	httpClient *http.Client
}

func New(maxConcurrentConnections int, allowInsecureTavern bool) *Fetcher {
	fetcher := &Fetcher{
		maxConns: maxConcurrentConnections,
	}
	if allowInsecureTavern {
		fetcher.httpClient = http3utils.DefaultClientInsecure
	} else {
		fetcher.httpClient = http3utils.DefaultClient
	}
	return fetcher
}

// Resolve fetches a list of serving peer from the tracker.
func (fetcher *Fetcher) Resolve(tracker *trackerv1.Tracker) ([]Peer, error) {
	// Attempt to collect all peers from all trackers
	peers := []Peer{}
	for _, tavern := range tracker.Taverns {
		tavernPeers, err := fetcher.ListPeers(tracker.Hash, tavern)
		if err != nil {
			return nil, fmt.Errorf("list peers: %v", err)
		}
		peers = append(peers, tavernPeers...)
	}
	// And update peer chunksets
	for i := range peers {
		peers[i].Refresh()
	}
	return peers, nil
}

// Fetch downloads a file to the given destination.
func (fetcher *Fetcher) Fetch(path string, tracker *trackerv1.Tracker) error {
	peers, err := fetcher.Resolve(tracker)
	if err != nil {
		return fmt.Errorf("resolve peers: %v", err)
	}
	// Spawn workers
	tasks := make(chan Task)
	for i := 0; i < fetcher.maxConns; i++ {
		go fetcher.handleFetchTasks(tasks)
	}
	// Track fetch progress
	fetched := make(chan int64, 1)
	done := make(chan bool)
	// Show progress in terminal
	progress := progressbar.DefaultBytes(tracker.Size, "Fetching chunks from peers")
	go func() {
		for !progress.IsFinished() {
			progress.Add64(<-fetched)
		}
		done <- true
	}()
	// Schedule chunk fetches
	for i := range tracker.ChunkHashes {
		tasks <- Task{
			Fetched:     fetched,
			Tracker:     tracker,
			Peers:       peers,
			Hash:        tracker.Hash,
			ChunkHash:   tracker.ChunkHashes[i],
			Chunk:       int64(i),
			Destination: generateChunkPath(path, i),
		}
	}
	// Synchronize, shut down workers
	<-done
	close(done)
	close(tasks)
	close(fetched)
	progress.Finish()
	// Stick chunks together
	progress = progressbar.DefaultBytes(tracker.Size, "Sticking chunks together")
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create final file: %v", err)
	}
	defer file.Close()
	for i := range tracker.ChunkHashes {
		if err := func() error {
			chunk, err := readAndVerify(generateChunkPath(path, i), tracker.ChunkHashes[i])
			if err != nil {
				return err
			}
			if _, err := file.Write(chunk); err != nil {
				return err
			}
			// Update progressbar
			progress.Add64(int64(len(chunk)))
			return nil
		}(); err != nil {
			return fmt.Errorf("final assembly: %v", err)
		}
	}
	// Cleanup chunk files
	for i := range tracker.ChunkHashes {
		chunkpath := fmt.Sprintf("%s.%d", path, i)
		if err := os.Remove(chunkpath); err != nil {
			return fmt.Errorf("remove chunkfile %s: %v", chunkpath, err)
		}
	}
	return nil
}
