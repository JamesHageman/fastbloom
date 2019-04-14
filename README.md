# FastBloom

[![CircleCI](https://circleci.com/gh/JamesHageman/fastbloom/tree/master.svg?style=svg)](https://circleci.com/gh/JamesHageman/fastbloom/tree/master)
[![codecov](https://codecov.io/gh/JamesHageman/fastbloom/branch/master/graph/badge.svg)](https://codecov.io/gh/JamesHageman/fastbloom)
[![GoDoc](https://godoc.org/github.com/JamesHageman/fastbloom?status.svg)](https://godoc.org/github.com/JamesHageman/fastbloom)
[![Go Report Card](https://goreportcard.com/badge/github.com/JamesHageman/fastbloom)](https://goreportcard.com/report/github.com/JamesHageman/fastbloom)

`fastbloom` implements a [Bloom Filter](https://en.wikipedia.org/wiki/Bloom_filter)
that supports any number of concurrent readers and writers.

Docs: http://godoc.org/github.com/JamesHageman/fastbloom

Example:

```go
package main

import "github.com/JamesHageman/fastbloom"

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
```

## Concurrency

`fastbloom` allows any number of concurrent readers and writers. It maintains
the following invariants:

- Before `Add(key)` or `TestAndAdd(key)` are called, `Test(key)` may return false.
- After `Add(key)` or `TestAndAdd(key)` are called, `Test(key)` will always return true.
- If `Test(key)` and either of the mutating methods (`Add` or `TestAndAdd`) are called
  concurrently, `Test` may return false.

Instead of relying on locks, `fastbloom` takes advantage of the `sync/atomic`
package from the standard library.

The bit vector is implemented as a `[]uint32`. Each `unint32` can be thought of
as a 32 bit block. When reading the bit vector, the reader queries bits by calling
`atomic.LoadUint32(addr)` on each block. When writing a bit, the writer uses a
combination of `atomic.LoadUint32(addr)` and `atomic.CompareAndSwapUint32(addr, orig, updated)`
to safely read a block, mutate it, and write it back to memory. If two concurrent writers
try to mutate the same block, both will only finish once `atomic.CompareAndSwapUint32` returns
true, meaning that the write-after read was guaranteed to be atomic. This ensures
that two writers that try to update different bits withing a block do not
unexpectedly erase the other's changes.
