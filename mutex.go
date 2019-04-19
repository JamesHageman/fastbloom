package fastbloom

import (
	"sync"
)

type MutexFilter struct {
	*LockFreeFilter
	sync.RWMutex
}

func NewMutexFilter(n uint, fpRate float64) *MutexFilter {
	return &MutexFilter{LockFreeFilter: NewFilter(n, fpRate)}
}

func (f *MutexFilter) Test(key []byte) bool {
	f.RLock()
	defer f.RUnlock()

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

func (f *MutexFilter) Add(key []byte) {
	f.Lock()
	defer f.Unlock()

	lower, upper := f.hash(key)

	// Set all k bits to 1
	for i := uint(0); i < f.k; i++ {
		offset := (uint(lower) + uint(upper)*i) % f.m
		f.setBit(offset)
	}
}

func (f *MutexFilter) TestAndAdd(key []byte) bool {
	f.Lock()
	defer f.Unlock()

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

func (f *MutexFilter) getBit(offset uint) bool {
	index := offset / 32
	bit := offset % 32
	mask := uint32(1 << bit)

	b := f.data[index]
	return b&mask != 0
}

func (f *MutexFilter) setBit(offset uint) {
	index := offset / 32
	bit := offset % 32
	mask := uint32(1 << bit)

	orig := f.data[index]
	updated := orig | mask
	f.data[index] = updated
}
