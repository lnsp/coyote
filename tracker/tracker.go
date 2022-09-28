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

package tracker

import (
	"crypto/sha256"
	fmt "fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	trackerv1 "github.com/lnsp/ftp2p/gen/tracker/v1"
	"google.golang.org/protobuf/proto"
)

// Open opens a tracker file, decodes it and returns its contents.
func Open(name string) (*trackerv1.Tracker, error) {
	trackfile, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("open tracker: %v", err)
	}
	var tracker trackerv1.Tracker
	if err := proto.Unmarshal(trackfile, &tracker); err != nil {
		return nil, fmt.Errorf("decode tracker: %v", err)
	}
	return &tracker, nil
}

// Save stores a tracker structure on disk.
func Save(tracker *trackerv1.Tracker, name string) error {
	trackfile, err := proto.Marshal(tracker)
	if err != nil {
		return fmt.Errorf("encode tracker: %v", err)
	}
	if err := ioutil.WriteFile(name, trackfile, 0644); err != nil {
		return fmt.Errorf("write tracker: %v", err)
	}
	return nil
}

const MinChunkCount int64 = 8
const MaxChunkSize int64 = 2 << 16

// Track generates a new tracker file that is ready for seeding.
func Track(addr string, path string) (*trackerv1.Tracker, error) {
	// Check if path is folder or does not exist or something
	fstat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat seedfile: %v", err)
	} else if fstat.IsDir() {
		return nil, fmt.Errorf("seedfile can not be directory")
	}
	// Generate tracker structure
	// Optimize chunk size, we want to have at least 8 chunks with max size 64KiB
	fileSize := fstat.Size()
	chunkSize := MaxChunkSize
	chunkCount := fileSize / chunkSize
	if chunkCount < MinChunkCount {
		chunkCount = MinChunkCount
		chunkSize = fileSize / chunkCount
	}
	// Scan through file, do chunk hashes
	fileHash, chunkHash := sha256.New(), sha256.New()
	hash := io.MultiWriter(fileHash, chunkHash)
	hashes := make([][]byte, chunkCount)
	seedfile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open seedfile: %v", err)
	}
	defer seedfile.Close()
	chunk := make([]byte, fileSize-(chunkCount-1)*chunkSize)
	for index := int64(0); index < chunkCount; index++ {
		size := chunkSize
		if index == chunkCount-1 {
			size = fileSize - (chunkCount-1)*chunkSize
		}
		_, err := seedfile.Read(chunk[:size])
		if err != nil {
			return nil, fmt.Errorf("scan seedfile: %v", err)
		}
		chunkHash.Reset()
		hash.Write(chunk[:size])
		hashes[index] = chunkHash.Sum(nil)
	}
	return &trackerv1.Tracker{
		Addr:        addr,
		Hash:        fileHash.Sum(nil),
		Name:        filepath.Base(path),
		Size:        fileSize,
		ChunkSize:   chunkSize,
		ChunkHashes: hashes,
	}, nil
}
