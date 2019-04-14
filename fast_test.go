package fastbloom

import (
	"strconv"
	"sync"
	"testing"
)

// Ensures that Capacity returns the number of bits, m, in the Bloom filter.
func TestLockFreeBloomFilter_Capacity(t *testing.T) {
	f := NewFilter(100, 0.1)

	if capacity := f.Capacity(); capacity != 480 {
		t.Errorf("Expected 480, got %d", capacity)
	}
}

// Ensures that K returns the number of hash functions in the Bloom Filter.
func TestLockFreeBloomFilter_K(t *testing.T) {
	f := NewFilter(100, 0.1)

	if k := f.K(); k != 4 {
		t.Errorf("Expected 4, got %d", k)
	}
}

// Ensures that Test, Add, and TestAndAdd behave correctly.
func TestLockFreeBloomFilter_TestAndAdd(t *testing.T) {
	f := NewFilter(100, 0.01)

	// `a` isn't in the filter.
	if f.Test([]byte(`a`)) {
		t.Error("`a` should not be a member")
	}

	if f.Add([]byte(`a`)) != f {
		t.Error("Returned BloomFilter should be the same instance")
	}

	// `a` is now in the filter.
	if !f.Test([]byte(`a`)) {
		t.Error("`a` should be a member")
	}

	// `a` is still in the filter.
	if !f.TestAndAdd([]byte(`a`)) {
		t.Error("`a` should be a member")
	}

	// `b` is not in the filter.
	if f.TestAndAdd([]byte(`b`)) {
		t.Error("`b` should not be a member")
	}

	// `a` is still in the filter.
	if !f.Test([]byte(`a`)) {
		t.Error("`a` should be a member")
	}

	// `b` is now in the filter.
	if !f.Test([]byte(`b`)) {
		t.Error("`b` should be a member")
	}

	// `c` is not in the filter.
	if f.Test([]byte(`c`)) {
		t.Error("`c` should not be a member")
	}

	for i := 0; i < 1000000; i++ {
		f.TestAndAdd([]byte(strconv.Itoa(i)))
	}

	// `x` should be a false positive.
	if !f.Test([]byte(`x`)) {
		t.Error("`x` should be a member")
	}
}

func TestLockFreeBloomFilter_Add_Concurrent(t *testing.T) {
	workers := 100
	perWorker := 1000
	n := uint(perWorker * workers)
	f := NewFilter(n, 0.01)

	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for j := 0; j < perWorker; j++ {
				k := start + j
				key := []byte(strconv.Itoa(k))
				f.Add(key)
			}
		}(perWorker * i)
	}
	wg.Wait()

	for i := 0; i < int(n); i++ {
		key := []byte(strconv.Itoa(i))
		if !f.Test(key) {
			t.Errorf("key `%s` should be a member", string(key))
		}
	}
}

func BenchmarkLockFreeAdd(b *testing.B) {
	b.StopTimer()
	f := NewFilter(uint(b.N), 0.1)
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		f.Add(data[n])
	}
}

func BenchmarkLockFreeAdd_Concurrent(b *testing.B) {
	b.StopTimer()
	f := NewFilter(uint(b.N), 0.1)
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	workers := 4
	wg := sync.WaitGroup{}
	wg.Add(workers)
	b.StartTimer()

	for w := 0; w < workers; w++ {
		start := b.N / workers * w
		end := b.N / workers * (w + 1)

		go func() {
			for i := start; i < end; i++ {
				f.Add(data[i])
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
