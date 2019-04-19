package fastbloom_test

import (
	"testing"

	"github.com/JamesHageman/fastbloom"
)

func BenchmarkMutexFilter_Add_SingleThreaded(b *testing.B) {
	b.StopTimer()
	f := fastbloom.NewMutexFilter(uint(b.N), fpRate)
	b.StartTimer()
	benchmarkFilterAdd_SingleThreaded(b, f.Add)
}

func BenchmarkMutexFilter_Add_4ConcurrentWriters(b *testing.B) {
	b.StopTimer()
	f := fastbloom.NewMutexFilter(uint(b.N), fpRate)
	b.StartTimer()
	benchmarkFilter_Add_4ConcurrentWriters(b, f.Add)
}

func BenchmarkMutexFreeFilter_Test_SingleThreaded(b *testing.B) {
	b.StopTimer()
	f := fastbloom.NewMutexFilter(uint(b.N), fpRate)
	b.StartTimer()
	benchmarkFilterTest_SingleThreaded(b, f.Add, f.Test)
}

func BenchmarkMutexFreeFilter_Test_4ConcurrentReaders(b *testing.B) {
	b.StopTimer()
	f := fastbloom.NewMutexFilter(uint(b.N), fpRate)
	b.StartTimer()
	benchmarkFilter_Test_4ConcurrentReaders(b, f.Add, f.Test)
}
