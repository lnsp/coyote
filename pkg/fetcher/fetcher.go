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
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/lnsp/ftp2p/pkg/seeder"
	"github.com/lnsp/ftp2p/pkg/tavern"
	"github.com/lnsp/ftp2p/pkg/tracker"
	"github.com/lnsp/grpc-quic/opts"

	qgrpc "github.com/lnsp/grpc-quic"
	"github.com/willf/bitset"
)

// ListPeers contacts a tracker looking for a given hash.
func ListPeers(hash []byte, addr string) ([]Peer, error) {
	// Contact tracker for peers
	tavernConn, err := qgrpc.Dial(addr, opts.WithInsecure(), opts.WithTimeout(time.Minute), opts.WithBackoffMaxDelay(time.Minute))
	if err != nil {
		return nil, fmt.Errorf("dial tavern: %v", err)
	}
	defer tavernConn.Close()
	tavernClient := tavern.NewTavernClient(tavernConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	resp, err := tavernClient.List(ctx, &tavern.ListRequest{
		Hash: hash,
	})
	if err != nil {
		return nil, fmt.Errorf("list peers: %v", err)
	}
	peers := make([]Peer, len(resp.Peers))
	for i := range resp.Peers {
		peers[i] = Peer{
			Address: resp.Peers[i].Addr,
			Hash:    hash,
		}
	}
	return peers, nil
}

// HasChunks fetches the set of chunks served by the peer.
func HasChunks(hash []byte, peer string) (*bitset.BitSet, error) {
	peerConn, err := qgrpc.Dial(peer, opts.WithInsecure(), opts.WithTimeout(time.Minute), opts.WithBackoffMaxDelay(time.Minute))
	if err != nil {
		return nil, fmt.Errorf("dial peer: %v", err)
	}
	defer peerConn.Close()
	peerClient := seeder.NewSeederClient(peerConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	resp, err := peerClient.Has(ctx, &seeder.HasRequest{
		Hash: hash,
	})
	if err != nil {
		return nil, fmt.Errorf("seeder has: %v", err)
	}
	return bitset.From(resp.Chunks), nil
}

type Peer struct {
	Address  string
	Hash     []byte
	Chunkset *bitset.BitSet
}

func (peer *Peer) Refresh() {
	chunkset, err := HasChunks(peer.Hash, peer.Address)
	if err != nil {
		log.Printf("fetch chunkset from peer %s: %v", peer.Address, err)
		peer.Chunkset = bitset.New(0)
	} else {
		peer.Chunkset = chunkset
	}
}

type Task struct {
	Fetched     chan int64
	Tracker     *tracker.Tracker
	Peers       []Peer
	Peer        int
	Hash        []byte
	ChunkHash   []byte
	Chunk       int64
	Destination string
}

func (task Task) Fetch() error {
	peerConn, err := qgrpc.Dial(task.Peers[task.Peer].Address, opts.WithInsecure(), opts.WithTimeout(time.Minute), opts.WithBackoffMaxDelay(time.Minute))
	if err != nil {
		return fmt.Errorf("dial peer: %v", err)
	}
	defer peerConn.Close()
	peerClient := seeder.NewSeederClient(peerConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	resp, err := peerClient.Fetch(ctx, &seeder.FetchRequest{
		Hash:  task.Hash,
		Chunk: task.Chunk,
	})
	if err != nil {
		return fmt.Errorf("seeder fetch: %v", err)
	}
	hasher := sha256.New()
	hasher.Write(resp.Data)
	if !bytes.Equal(hasher.Sum(nil), task.ChunkHash) {
		return fmt.Errorf("bad chunk from peer")
	}
	if err := ioutil.WriteFile(task.Destination, resp.Data, 0644); err != nil {
		return fmt.Errorf("chunk write: %v", err)
	}
	return nil
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

func fetchWorker(tasks chan Task) {
	for task := range tasks {
		chunk := uint(task.Chunk)
		fetched := false
		// Check for local correct copy
		if _, err := readAndVerify(task.Destination, task.ChunkHash); err == nil {
			fetched = true
		}
		// If not exist, try to fetch
		for cycle := 0; !fetched; cycle++ {
			// To load balance, generate random peer order per cycle
			order := randomIntOrder(len(task.Peers))
			for i := 0; i < len(task.Peers) && !fetched; i++ {
				j := order[i]
				if !task.Peers[j].Chunkset.Test(chunk) {
					log.Printf("peer %s does not have chunk %d", task.Peers[j].Address, chunk)
					continue
				}
				task.Peer = j
				if err := task.Fetch(); err != nil {
					log.Printf("fetch chunk %d from peer %s failed: %v", chunk, task.Peers[j].Address, err)
					continue
				}
				fetched = true
			}
			if !fetched {
				log.Printf("worker cycled for chunk %d in iteration %d", chunk, cycle)
				// Do exponential backoff, reaches max around 16 cycles
				backoff := time.Second * time.Duration(math.Min(60.0, math.Pow(1.3, float64(cycle))))
				time.Sleep(backoff)
				// To resolve worker cycle, re-fetch peers and chunksets
				peers, err := Resolve(task.Tracker)
				if err != nil {
					log.Printf("resolve after cycle: %v", err)
				} else {
					task.Peers = peers
				}
			}
		}
		// Submit success
		task.Fetched <- task.Chunk
	}
}

func readAndVerify(path string, hash []byte) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
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

// Resolve fetches a list of serving peer from the tracker.
func Resolve(tracker *tracker.Tracker) ([]Peer, error) {
	peers, err := ListPeers(tracker.Hash, tracker.Addr)
	if err != nil {
		return nil, fmt.Errorf("list peers: %v", err)
	}
	// And update peer chunksets
	for i := range peers {
		peers[i].Refresh()
	}
	return peers, nil
}

// Fetch downloads a file to the given destination.
func Fetch(path string, tracker *tracker.Tracker, numWorkers int) error {
	peers, err := Resolve(tracker)
	if err != nil {
		return fmt.Errorf("resolve peers: %v", err)
	}
	// Spawn workers
	tasks := make(chan Task)
	for i := 0; i < numWorkers; i++ {
		go fetchWorker(tasks)
	}
	// Track fetch progress
	numFetched, numChunks := 0, len(tracker.ChunkHashes)
	fetched := make(chan int64, 1)
	done := make(chan bool)
	go func() {
		for numFetched < numChunks {
			<-fetched
			numFetched++
			log.Printf("total progress: %.2f%%", float32(numFetched)/float32(numChunks)*100.)
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
	// Stick chunks together
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
