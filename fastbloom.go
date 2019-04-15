package fastbloom

import (
	"hash"
	"hash/fnv"
	"math"
	"sync"
	"sync/atomic"
)

// Filter is an implementation of a Bloom Filter that permits n concurrent
// readers and m concurrent writers.
type Filter struct {
	data []uint32 // underlying bit vector
	m    uint     // filter size
	k    uint     // number of hash functions

	hashPool sync.Pool // pool of hash objects
}

const fillRatio = 0.5

// NewFilter creates a bloom filter optimized for n elements with a falsePositive rate fpRate.
func NewFilter(n uint, fpRate float64) *Filter {
	m := optimalM(n, fpRate)
	return &Filter{
		data: make([]uint32, m/32+1),
		m:    m,
		k:    optimalK(fpRate),
		hashPool: sync.Pool{
			New: func() interface{} {
				return fnv.New64()
			},
		},
	}
}

// Capacity returns the Bloom filter capacity, m.
func (f *Filter) Capacity() uint {
	return f.m
}

// K returns the number of hash functions.
func (f *Filter) K() uint {
	return f.k
}

// Test tests the filter for the presence of a key. It is (only) guaranteed to return
// true after a call to Add(key) or TestAndAdd(key) have completed. Because of
// the possibility of false positives, Test(key) could also return true if the key
// hasn't been added.
func (f *Filter) Test(key []byte) bool {
	lower, upper := f.hash(key)

	// If any of the K bits are not set, then it's not a member.
	for i := uint(0); i < f.k; i++ {
		offset := (uint(lower) + uint(upper)*i) % f.m

		if !f.getBit(offset) {
			return false
		}
	}

	return true
}

// Add writes a key to the bloom filter. It can be called concurrently with other
// calls to Add, TestAndAdd, and Test.
func (f *Filter) Add(key []byte) *Filter {
	lower, upper := f.hash(key)

	// Set all k bits to 1
	for i := uint(0); i < f.k; i++ {
		offset := (uint(lower) + uint(upper)*i) % f.m
		f.setBit(offset)
	}

	return f
}

// TestAndAdd adds a key to the bloom filter and returns true if it appears that
// the key was already present.
func (f *Filter) TestAndAdd(key []byte) bool {
	lower, upper := f.hash(key)
	member := true

	// If any of the K bits are not set, then it's not a member.
	for i := uint(0); i < f.k; i++ {
		offset := (uint(lower) + uint(upper)*i) % f.m

		if !f.getBit(offset) {
			member = false
			f.setBit(offset)
		}
	}

	return member
}

func (f *Filter) hash(data []byte) (uint32, uint32) {
	h := f.hashPool.Get().(hash.Hash64)
	h.Write(data)
	sum := h.Sum64()
	higher := uint32(sum >> 32)
	lower := uint32(sum)
	h.Reset()
	f.hashPool.Put(h)
	return higher, lower
}

func (f *Filter) getBit(offset uint) bool {
	byteIndex := offset / 32
	byteOffset := offset % 32
	mask := uint32(1 << byteOffset)

	b := atomic.LoadUint32(&f.data[byteIndex])
	return b&mask != 0
}

func (f *Filter) setBit(offset uint) {
	index := offset / 32
	bit := offset % 32
	mask := uint32(1 << bit)

	ptr := &f.data[index]

	for {
		orig := atomic.LoadUint32(ptr)
		updated := orig | mask
		swapped := atomic.CompareAndSwapUint32(ptr, orig, updated)
		if swapped {
			break
		}
	}
}

// optimalM calculates the optimal Bloom filter size, m, based on the number of
// items and the desired rate of false positives.
func optimalM(n uint, fpRate float64) uint {
	return uint(math.Ceil(float64(n) / ((math.Log(fillRatio) *
		math.Log(1-fillRatio)) / math.Abs(math.Log(fpRate)))))
}

// optimalK calculates the optimal number of hash functions to use for a Bloom
// filter based on the desired rate of false positives.
func optimalK(fpRate float64) uint {
	return uint(math.Ceil(math.Log2(1 / fpRate)))
}
