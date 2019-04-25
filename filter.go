package fastbloom

import (
	"hash"
	"hash/fnv"
	"sync"
)

type BloomFilter struct {
	m        uint
	k        uint
	buckets  Buckets
	hash     func() hash.Hash64
	hashPool sync.Pool
}

func NewDefaultBloomFilter(n uint, fpRate float64) *BloomFilter {
	return NewBloomFilter(optimalM(n, fpRate), optimalK(fpRate))
}

func NewBloomFilter(m, k uint, config ...func(filter *BloomFilter)) *BloomFilter {
	filter := &BloomFilter{
		m: m,
		k: k,
	}

	for _, f := range config {
		f(filter)
	}

	if filter.hash == nil {
		filter.hash = fnv.New64
	}

	filter.hashPool = sync.Pool{
		New: func() interface{} {
			return filter.hash()
		},
	}

	if filter.buckets == nil {
		filter.buckets = newBuckets32(filter.m)
	}

	return filter
}

func (f *BloomFilter) Add(key []byte) {
	lower, upper := f.computeHash(key)

	// Set all k bits to 1
	for i := uint(0); i < f.k; i++ {
		offset := (uint(lower) + uint(upper)*i) % f.m
		f.buckets.SetBit(offset)
	}
}

func (f *BloomFilter) AddString(key string) {
	f.Add([]byte(key))
}

func (f *BloomFilter) Test(key []byte) bool {
	lower, upper := f.computeHash(key)

	// If any of the k bits are not set, then key not a member.
	for i := uint(0); i < f.k; i++ {
		offset := (uint(lower) + uint(upper)*i) % f.m

		if !f.buckets.GetBit(offset) {
			return false
		}
	}

	return true
}

func (f *BloomFilter) TestString(key string) bool {
	return f.Test([]byte(key))
}

func (f *BloomFilter) computeHash(key []byte) (higher uint32, lower uint32) {
	h := f.hashPool.Get().(hash.Hash64)
	h.Write(key)
	sum := h.Sum64()
	higher = uint32(sum >> 32)
	lower = uint32(sum)
	h.Reset()
	f.hashPool.Put(h)
	return
}
