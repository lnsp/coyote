package tavern

import (
	"bytes"
	"sort"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	tavernv1 "github.com/lnsp/coyote/gen/tavern/v1"
)

type PeerEntry struct {
	Peer      *tavernv1.Peer
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

func (index *PeerIndex) Seed(hash []byte, peer *tavernv1.Peer) {
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
