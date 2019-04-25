package fastbloom_test

import (
	"strconv"
	"sync"
	"testing"
)

const fpRate float64 = 0.01

func benchmarkFilterAdd_SingleThreaded(b *testing.B, add func(key []byte)) {
	b.StopTimer()
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		add(data[n])
	}
}

func benchmarkFilter_Add_4ConcurrentWriters(b *testing.B, add func(key []byte)) {
	b.StopTimer()
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
				add(data[i])
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func benchmarkFilterTest_SingleThreaded(b *testing.B, add func(key []byte), test func(key []byte) bool) {
	b.StopTimer()
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
		add(data[i])
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		test(data[n])
	}
}

func benchmarkFilter_Test_4ConcurrentReaders(b *testing.B, add func(key []byte), test func(key []byte) bool) {
	b.StopTimer()
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
		add(data[i])
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
				test(data[i])
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
