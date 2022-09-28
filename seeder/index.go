package seeder

import (
	"bytes"
	"sort"
	"sync"

	"github.com/bits-and-blooms/bitset"
	"github.com/emirpasic/gods/maps/treemap"
	trackerv1 "github.com/lnsp/ftp2p/gen/tracker/v1"
	"github.com/lnsp/ftp2p/tracker"
)

func byteSliceComparator(a, b interface{}) int {
	return bytes.Compare(a.([]byte), b.([]byte))
}

type FileEntry struct {
	Path   string
	Chunk  int64
	Offset int64
	Size   int64
}

type FileIndex struct {
	sync.RWMutex
	treemap *treemap.Map
}

func (index *FileIndex) Scan(hash []byte) *bitset.BitSet {
	index.RLock()
	value, ok := index.treemap.Get(hash)
	index.RUnlock()
	if !ok {
		return bitset.New(0)
	}
	chunks := bitset.New(uint(tracker.MinChunkCount))
	entries := value.([]*FileEntry)
	for _, e := range entries {
		chunks.Set(uint(e.Chunk))
	}
	return chunks
}

func (index *FileIndex) Get(hash []byte, chunk int64) (*FileEntry, bool) {
	index.RLock()
	value, ok := index.treemap.Get(hash)
	index.RUnlock()
	if !ok {
		return nil, false
	}
	entries := value.([]*FileEntry)
	found := sort.Search(len(entries), func(i int) bool {
		return entries[i].Chunk >= chunk
	})
	if found >= len(entries) || entries[found].Chunk != chunk {
		return nil, false
	}
	return entries[found], true
}

func (index *FileIndex) Add(path string, tracker *trackerv1.Tracker) {
	var offset int64
	entries := make([]*FileEntry, len(tracker.ChunkHashes))
	for chunk := range tracker.ChunkHashes {
		size := tracker.ChunkSize
		if chunk == len(tracker.ChunkHashes)-1 {
			size = tracker.Size - offset
		}
		entries[chunk] = &FileEntry{
			Path:   path,
			Chunk:  int64(chunk),
			Offset: offset,
			Size:   size,
		}
		offset += size
	}
	index.Lock()
	index.treemap.Put(tracker.Hash, entries)
	index.Unlock()
}

func newFileIndex() *FileIndex {
	return &FileIndex{
		treemap: treemap.NewWith(byteSliceComparator),
	}
}
