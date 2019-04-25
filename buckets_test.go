package fastbloom

import (
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testBuckets(t *testing.T, bucket func(m uint) Buckets) {
	t.Run("single threaded", func(t *testing.T) {
		m := uint(100)
		b := bucket(m)
		for i := uint(0); i < m; i++ {
			assert.Equal(t, false, b.GetBit(i), "GetBit(%d) should return false by default", i)
		}

		for i := uint(0); i < m; i++ {
			if i > 0 {
				assert.Equal(t, true, b.GetBit(i-1), "GetBit(%d) should return true when set", i)
			}

			b.SetBit(i)
			assert.Equal(t, true, b.GetBit(i), "GetBit(%d) should return true when set", i)

			if i < m-1 {
				assert.Equal(t, false, b.GetBit(i+1), "GetBit(%d) should return false when not set", i)
			}
		}
	})

	t.Run("concurrent reads and writes", func(t *testing.T) {
		m := uint(1000)
		b := bucket(m)

		var wg sync.WaitGroup
		for i := uint(0); i < m; i++ {
			i := i

			// spin off a reader routine to potentially detect races with the race detector
			wg.Add(1)
			go func() { defer wg.Done(); b.GetBit(i) }()

			// spin off a writer routine to write every individual bit
			wg.Add(1)
			go func() { defer wg.Done(); b.SetBit(i) }()
		}
		wg.Wait()

		// check that all bits were successfully set
		for i := uint(0); i < m; i++ {
			assert.Equal(t, true, b.GetBit(i), "GetBit(%d) should return true", i)
		}
	})

	t.Run("encode decode", func(t *testing.T) {
		m := uint(1000)
		b := bucket(m)
		for i := uint(0); i < m; i++ {
			if i%3 == 0 {
				b.SetBit(i)
			}
		}

		buf, err := b.GobEncode()
		assert.NoError(t, err)

		decoded := bucket(0)
		err = decoded.GobDecode(buf)
		assert.NoError(t, err)

		for i := uint(0); i < m; i++ {
			assert.Equal(t, i%3 == 0, b.GetBit(i))
		}
	})

	t.Run("concurrent encode", func(t *testing.T) {
		m := uint(1000)
		b := bucket(m)
		b.SetBit(0) // set one bit to 1

		var wg sync.WaitGroup
		var buf []byte

		for i := uint(0); i < m; i++ {
			i := i

			// spin off a writer routine to write every individual bit
			wg.Add(1)
			go func() { defer wg.Done(); b.SetBit(i) }()
		}

		// encode the bucket while being concurrently
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			buf, err = b.GobEncode()
			assert.NoError(t, err)
		}()

		wg.Wait()

		// every bit should be set...
		for i := uint(0); i < m; i++ {
			assert.Equal(t, true, b.GetBit(i))
		}

		// ...even if we clear the encoded buffer.
		for i := range buf {
			buf[i] = 0
		}
		for i := uint(0); i < m; i++ {
			assert.Equal(t, true, b.GetBit(i))
		}
	})
}

func benchmarkBuckets(b *testing.B, bucket func(m uint) Buckets) {
	b.Run("single threaded reads", func(b *testing.B) {
		b.StopTimer()
		bs := bucket(uint(b.N))
		for i := 0; i < b.N; i++ {
			if i%2 == 0 {
				bs.SetBit(uint(i))
			}
		}
		b.StartTimer()
		for i, n := uint(0), uint(b.N); i < n; i++ {
			bs.GetBit(i)
		}
	})

	b.Run("single threaded writes", func(b *testing.B) {
		b.StopTimer()
		bs := bucket(uint(b.N))
		for i := 0; i < b.N; i++ {
			if i%2 == 0 {
				bs.SetBit(uint(i))
			}
		}
		b.StartTimer()
		for i, n := uint(0), uint(b.N); i < n; i++ {
			bs.SetBit(i)
		}
	})

	workers := uint(runtime.GOMAXPROCS(0))

	b.Run("concurrent readers", func(b *testing.B) {
		b.StopTimer()
		m := uint(b.N)
		bs := bucket(m)
		for i := 0; i < b.N; i++ {
			if i%2 == 0 {
				bs.SetBit(uint(i))
			}
		}
		var wg sync.WaitGroup
		wg.Add(int(workers))

		b.StartTimer()
		for i := uint(0); i < workers; i++ {
			i := i
			go func() {
				defer wg.Done()
				for j := i; j < m; j += workers {
					bs.GetBit(j)
				}
			}()
		}
		wg.Wait()
	})
}

func TestBuckets32(t *testing.T) {
	testBuckets(t, func(m uint) Buckets {
		return newBuckets32(m)
	})
}

func BenchmarkBuckets32(b *testing.B) {
	benchmarkBuckets(b, func(m uint) Buckets {
		return newBuckets32(m)
	})
}

func TestBuckets64(t *testing.T) {
	testBuckets(t, func(m uint) Buckets {
		return newBuckets64(m)
	})
}

func BenchmarkBuckets64(b *testing.B) {
	benchmarkBuckets(b, func(m uint) Buckets {
		return newBuckets64(m)
	})
}

func TestMutexBuckets(t *testing.T) {
	testBuckets(t, func(m uint) Buckets {
		return newMutexBuckets(m)
	})
}

func BenchmarkMutexBuckets(b *testing.B) {
	benchmarkBuckets(b, func(m uint) Buckets {
		return newMutexBuckets(m)
	})
}
