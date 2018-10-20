package lsh

import (
	"math"
)

// Hash is a hash function.
type Hash func(int, int) int

var (
	hashes = []Hash{
		func(x, numBuckets int) int {
			return (x + 1) % numBuckets
		},
		func(x, numBuckets int) int {
			return (3*x + 1) % numBuckets
		},
		func(x, numBuckets int) int {
			return (2*x + 4) % numBuckets
		},
		func(x, numBuckets int) int {
			return (3*x - 1) % numBuckets
		},
		func(x, numBuckets int) int {
			return (2*x + 1) % numBuckets
		},
		func(x, numBuckets int) int {
			return (3*x + 2) % numBuckets
		},
		func(x, numBuckets int) int {
			return (5*x + 2) % numBuckets
		},
		func(x, numBuckets int) int {
			return x % numBuckets
		},
		func(x, numBuckets int) int {
			return (31*x + x&0xff) % numBuckets
		},
		func(x, numBuckets int) int {
			var i = (x * 31) >> 28
			return (i & 15) % numBuckets
		},
		func(x, numBuckets int) int {
			return (5*x + 11) % numBuckets
		},
		func(x, numBuckets int) int {
			return x & math.MaxInt32 % numBuckets
		},
	}
)
