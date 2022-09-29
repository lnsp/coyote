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
	"fmt"
	"io"
	"os"
	"path/filepath"

	trackerv1 "github.com/lnsp/ftp2p/gen/tracker/v1"
	"github.com/schollz/progressbar"
	"google.golang.org/protobuf/proto"
)

// Open opens a tracker file, decodes it and returns its contents.
func Open(name string) (*trackerv1.Tracker, error) {
	trackfile, err := os.ReadFile(name)
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
	if err := os.WriteFile(name, trackfile, 0644); err != nil {
		return fmt.Errorf("write tracker: %v", err)
	}
	return nil
}

const minChunkCount int64 = 8
const maxChunkSize int64 = 1 << 24

// Track generates a new tracker file that is ready for seeding.
func Track(taverns []string, path string) (*trackerv1.Tracker, error) {
	// Check if path is folder or does not exist or something
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat seedfile: %v", err)
	} else if stat.IsDir() {
		return nil, fmt.Errorf("seedfile can not be directory")
	}
	// Generate tracker structure
	// Optimize chunk size, we want to have at least 8 chunks with max size 64KiB
	filesize := stat.Size()
	chunksize := maxChunkSize
	chunkcount := filesize / chunksize
	// Print out track chunksize and chunkcount
	if chunkcount < minChunkCount {
		chunkcount = minChunkCount
		chunksize = filesize / chunkcount
	}
	// Scan through file, do chunk hashes
	filehash, chunkhash := sha256.New(), sha256.New()
	hashwriter := io.MultiWriter(filehash, chunkhash)
	chunkhashes := make([][]byte, chunkcount)
	seedfile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open seedfile: %v", err)
	}
	defer seedfile.Close()
	chunk := make([]byte, filesize-(chunkcount-1)*chunksize)
	// Start progress bar for nice visuals
	pb := progressbar.New(int(chunkcount))
	for index := int64(0); index < chunkcount; index++ {
		size := chunksize
		if index == chunkcount-1 {
			size = filesize - (chunkcount-1)*chunksize
		}
		_, err := seedfile.Read(chunk[:size])
		if err != nil {
			return nil, fmt.Errorf("scan seedfile: %v", err)
		}
		chunkhash.Reset()
		hashwriter.Write(chunk[:size])
		chunkhashes[index] = chunkhash.Sum(nil)

		// Update progress bar
		pb.Add(1)
	}
	return &trackerv1.Tracker{
		Taverns:     taverns,
		Hash:        filehash.Sum(nil),
		Name:        filepath.Base(path),
		Size:        filesize,
		ChunkSize:   chunksize,
		ChunkHashes: chunkhashes,
	}, nil
}
