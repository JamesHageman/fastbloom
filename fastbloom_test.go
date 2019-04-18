package fastbloom_test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/JamesHageman/fastbloom"
	"github.com/stretchr/testify/assert"
)

// Ensures that Capacity returns the number of bits, m, in the Bloom filter.
func TestLockFreeBloomFilter_Capacity(t *testing.T) {
	f := fastbloom.NewFilter(100, 0.1)

	if capacity := f.Capacity(); capacity != 480 {
		t.Errorf("Expected 480, got %d", capacity)
	}
}

// Ensures that K returns the number of hash functions in the Bloom Filter.
func TestLockFreeBloomFilter_K(t *testing.T) {
	f := fastbloom.NewFilter(100, 0.1)

	if k := f.K(); k != 4 {
		t.Errorf("Expected 4, got %d", k)
	}
}

// Ensures that Test, Add, and TestAndAdd behave correctly.
func TestLockFreeBloomFilter_TestAndAdd(t *testing.T) {
	f := fastbloom.NewFilter(100, 0.01)

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
	f := fastbloom.NewFilter(n, 0.01)

	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		// create a worker to write data to the bloom filter
		go func(start int) {
			defer wg.Done()

			for j := 0; j < perWorker; j++ {
				k := start + j
				key := []byte(strconv.Itoa(k))
				f.Add(key)
			}
		}(perWorker * i)

		wg.Add(1)
		// create another worker to read data from the filter, which should trigger
		// the race detector if there is a data race.
		go func(start int) {
			defer wg.Done()

			for j := 0; j < perWorker; j++ {
				k := start + j
				key := []byte(strconv.Itoa(k))
				f.Test(key)
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

func TestFilter_GobEncode(t *testing.T) {
	keys := []string{`a`, `b`, `c`}
	filter := fastbloom.NewFilter(100, 0.01)
	for _, k := range keys {
		filter.Add([]byte(k))
	}
	b, err := filter.GobEncode()
	assert.NoError(t, err)
	decoded := &fastbloom.Filter{}
	err = decoded.GobDecode(b)
	assert.NoError(t, err)
	assert.Equal(t, filter.Capacity(), decoded.Capacity())
	assert.Equal(t, filter.K(), decoded.K())

	for _, k := range keys {
		assert.True(t, filter.Test([]byte(k)), "`%s` should be a member", k)
	}
}

func TestFilter_GobDecode(t *testing.T) {
	bs := []byte("foo")
	f := fastbloom.Filter{}
	err := f.GobDecode(bs)
	assert.Error(t, err)
}

func BenchmarkLockFreeAdd_SingleThreaded(b *testing.B) {
	b.StopTimer()
	f := fastbloom.NewFilter(uint(b.N), 0.1)
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		f.Add(data[n])
	}
}

func BenchmarkLockFreeAdd_4ConcurrentWriters(b *testing.B) {
	b.StopTimer()
	f := fastbloom.NewFilter(uint(b.N), 0.1)
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

func ExampleNewFilter() {
	filter := fastbloom.NewFilter(100, 0.01)
	fmt.Println(filter.Test([]byte(`a`)))
	fmt.Println(filter.Test([]byte(`b`)))

	filter.Add([]byte(`a`))
	filter.Add([]byte(`b`))

	fmt.Println(filter.Test([]byte(`a`)))
	fmt.Println(filter.Test([]byte(`b`)))
	fmt.Println(filter.Test([]byte(`c`)))

	// Output:
	// false
	// false
	// true
	// true
	// false
}
