package lsh

import (
	"math"
)

// Hash is a hash function.
type Hash func(float64) float64

var (
	hashes = []Hash{
		func(x float64) float64 {
			return math.Mod(x+1, 5)
		},
		func(x float64) float64 {
			return math.Mod(3*x+1, 5)
		},
		func(x float64) float64 {
			return math.Mod(x, 15)
		},
		func(x float64) float64 {
			var u32 = uint32(x*31) >> 28
			return float64(u32 & 15)
		},
	}
)
