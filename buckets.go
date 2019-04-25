package fastbloom

import (
	"encoding/binary"
	"encoding/gob"
	"sync"
	"sync/atomic"
)

type Buckets interface {
	SetBit(offset uint)
	GetBit(offset uint) bool
	gob.GobEncoder
	gob.GobDecoder
}

type buckets32 struct {
	data []uint32
}

func newBuckets32(m uint) *buckets32 {
	return &buckets32{data: make([]uint32, m/32+1)}
}

var _ Buckets = &buckets32{}

func (b *buckets32) GobEncode() ([]byte, error) {
	ret := make([]byte, 0, len(b.data)*4)
	buf := [4]byte{}
	for i := range b.data {
		block := atomic.LoadUint32(&b.data[i])
		binary.BigEndian.PutUint32(buf[:], block)
		ret = append(ret, buf[:]...)
	}
	return ret, nil
}

func (b *buckets32) GobDecode(data []byte) error {
	b.data = make([]uint32, len(data)/4+1)
	buf := [4]byte{}
	for i := 0; i < len(data); i += 4 {
		copy(buf[:], data[i:])
		block := binary.BigEndian.Uint32(buf[:])
		b.data[i/4] = block
	}
	return nil
}

func (b *buckets32) SetBit(offset uint) {
	index := offset / 32
	bit := offset % 32
	mask := uint32(1 << bit)
	ptr := &b.data[index]

	for {
		orig := atomic.LoadUint32(ptr)
		updated := orig | mask
		swapped := atomic.CompareAndSwapUint32(ptr, orig, updated)
		if swapped {
			break
		}
	}
}

func (b *buckets32) GetBit(offset uint) bool {
	index := offset / 32
	bit := offset % 32
	mask := uint32(1 << bit)
	ptr := &b.data[index]

	block := atomic.LoadUint32(ptr)
	return block&mask != 0
}

type buckets64 struct {
	data []uint64
}

var _ Buckets = &buckets64{}

func newBuckets64(m uint) *buckets64 {
	return &buckets64{data: make([]uint64, m/64+1)}
}

func (b *buckets64) SetBit(offset uint) {
	index := offset / 64
	bit := offset % 64
	mask := uint64(1 << bit)
	ptr := &b.data[index]

	for {
		orig := atomic.LoadUint64(ptr)
		updated := orig | mask
		swapped := atomic.CompareAndSwapUint64(ptr, orig, updated)
		if swapped {
			break
		}
	}
}

func (b *buckets64) GetBit(offset uint) bool {
	index := offset / 64
	bit := offset % 64
	mask := uint64(1 << bit)
	ptr := &b.data[index]

	block := atomic.LoadUint64(ptr)
	return block&mask != 0
}

func (b *buckets64) GobEncode() ([]byte, error) {
	ret := make([]byte, 0, len(b.data)*8)
	buf := [8]byte{}
	for i := range b.data {
		block := atomic.LoadUint64(&b.data[i])
		binary.BigEndian.PutUint64(buf[:], block)
		ret = append(ret, buf[:]...)
	}
	return ret, nil
}

func (b *buckets64) GobDecode(data []byte) error {
	b.data = make([]uint64, len(data)/8+1)
	buf := [8]byte{}
	for i := 0; i < len(data); i += 4 {
		copy(buf[:], data[i:])
		block := binary.BigEndian.Uint64(buf[:])
		b.data[i/8] = block
	}
	return nil
}

var _ Buckets = &mutexBuckets{}

func newMutexBuckets(m uint) *mutexBuckets {
	return &mutexBuckets{
		data: make([]byte, m/8+1),
	}
}

type mutexBuckets struct {
	lock sync.RWMutex
	data []byte
}

func (b *mutexBuckets) SetBit(offset uint) {
	index := offset / 8
	bit := offset % 8
	mask := byte(1 << bit)

	b.lock.Lock()
	defer b.lock.Unlock()

	orig := b.data[index]
	updated := orig | mask
	b.data[index] = updated
}

func (b *mutexBuckets) GetBit(offset uint) bool {
	index := offset / 8
	bit := offset % 8
	mask := byte(1 << bit)

	b.lock.RLock()
	defer b.lock.RUnlock()

	block := b.data[index]

	return block&mask != 0
}

func (b *mutexBuckets) GobEncode() ([]byte, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	buf := make([]byte, len(b.data))
	_ = copy(buf, b.data)
	return buf, nil
}

func (b *mutexBuckets) GobDecode(buf []byte) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.data = make([]byte, len(buf))
	copy(b.data, buf)
	return nil
}
