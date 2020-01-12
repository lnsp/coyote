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
	"github.com/lnsp/ftp2p/pkg/seeder"
	"github.com/lnsp/ftp2p/pkg/tavern"
	"github.com/lnsp/ftp2p/pkg/tracker"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/willf/bitset"
	"google.golang.org/grpc"
)

func listPeers(hash []byte, addr string) ([]string, error) {
	// Contact tracker for peers
	tavernConn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(time.Minute), grpc.WithBackoffMaxDelay(time.Minute))
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
	peers := make([]string, len(resp.Peers))
	for i := range resp.Peers {
		peers[i] = resp.Peers[i].Addr
	}
	return peers, nil
}

func hasChunks(hash []byte, peer string) (*bitset.BitSet, error) {
	peerConn, err := grpc.Dial(peer, grpc.WithInsecure(), grpc.WithTimeout(time.Minute), grpc.WithBackoffMaxDelay(time.Minute))
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

func fetchChunk(hash []byte, chunkhash []byte, peer string, chunk int64, path string) error {
	peerConn, err := grpc.Dial(peer, grpc.WithInsecure(), grpc.WithTimeout(time.Minute), grpc.WithBackoffMaxDelay(time.Minute))
	if err != nil {
		return fmt.Errorf("dial peer: %v", err)
	}
	defer peerConn.Close()
	peerClient := seeder.NewSeederClient(peerConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	resp, err := peerClient.Fetch(ctx, &seeder.FetchRequest{
		Hash:  hash,
		Chunk: chunk,
	})
	if err != nil {
		return fmt.Errorf("seeder fetch: %v", err)
	}
	hasher := sha256.New()
	hasher.Write(resp.Data)
	if !bytes.Equal(hasher.Sum(nil), chunkhash) {
		return fmt.Errorf("bad chunk from peer")
	}
	if err := ioutil.WriteFile(path, resp.Data, 0644); err != nil {
		return fmt.Errorf("chunk write: %v", err)
	}
	return nil
}

func Fetch(path string, tracker *tracker.Tracker) error {
	// List peers from tracker
	peers, err := listPeers(tracker.Hash, tracker.Addr)
	if err != nil {
		return fmt.Errorf("list peers: %v", err)
	}
	// Contact peers for chunks
	chunksets := make([]*bitset.BitSet, len(peers))
	for i := range peers {
		chunkset, err := hasChunks(tracker.Hash, peers[i])
		if err != nil {
			chunksets[i] = bitset.New(0)
		} else {
			chunksets[i] = chunkset
		}
	}
	// Look for chunks
	numChunks := uint(len(tracker.ChunkHashes))
	fetched := bitset.New(numChunks)
	for fetched.Count() < numChunks {
		for i, chunkHash := range tracker.ChunkHashes {
			path := fmt.Sprintf("%s.%d", path, i)
			for j, peer := range peers {
				if !chunksets[j].Test(uint(i)) {
					log.Printf("Peer %s does not have chunk %d", peer, i)
					continue
				}
				if err := fetchChunk(tracker.Hash, chunkHash, peer, int64(i), path); err != nil {
					log.Printf("Fetch chunk %d from peer %s failed: %v", i, peer, err)
					continue
				}
				fetched.Set(uint(i))
			}
		}
		progress := float32(fetched.Count()) / float32(numChunks) * 100.
		log.Printf("Total progress: %.2f%%", progress)
	}
	// Stick chunks together
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create final file: %v", err)
	}
	defer file.Close()
	for i := 0; i < len(tracker.ChunkHashes); i++ {
		if err := func() error {
			chunkpath := fmt.Sprintf("%s.%d", path, i)
			chunkfile, err := os.Open(chunkpath)
			if err != nil {
				return err
			}
			defer chunkfile.Close()
			if _, err := io.Copy(file, chunkfile); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return fmt.Errorf("final assembly: %v", err)
		}
	}
	return nil
}
