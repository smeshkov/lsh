package lsh

import (
	"fmt"
	"math"
	"math/rand"
)

// HashFunc is a hash function.
type HashFunc func(int, int) int

// Hasher is a holder of the specific hash function.
type Hasher struct {
	// related hash function
	hf HashFunc

	// string representation of the function
	t string
}

// Hash retrurns hash function of the hasher.
func (h *Hasher) Hash() HashFunc {
	return h.hf
}

func (h *Hasher) String() string {
	return h.t
}

// SuggestHashNum suggests number of generated hashes based on the average number of shingles.
func SuggestHashNum(avgNumOfShingles int) int {
	if avgNumOfShingles <= 100 {
		return 5
	} else if avgNumOfShingles > 100 && avgNumOfShingles <= 150 {
		return 10
	} else if avgNumOfShingles > 150 && avgNumOfShingles <= 300 {
		return 20
	} else if avgNumOfShingles > 300 && avgNumOfShingles <= 400 {
		return 30
	} else if avgNumOfShingles > 400 && avgNumOfShingles <= 550 {
		return 40
	} else if avgNumOfShingles > 550 && avgNumOfShingles <= 700 {
		return 50
	} else if avgNumOfShingles > 700 && avgNumOfShingles <= 800 {
		return 60
	} else if avgNumOfShingles > 800 && avgNumOfShingles <= 900 {
		return 70
	} else if avgNumOfShingles > 900 && avgNumOfShingles <= 1000 {
		return 80
	} else if avgNumOfShingles > 1000 && avgNumOfShingles <= 1200 {
		return 90
	}

	return 100
}

// GenerateHashers generates hash functions for given amount.
func GenerateHashers(amount int) []*Hasher {
	hashers := make([]*Hasher, amount)
	seen := make(map[string]bool)

	// simple modulus func is 1st
	modulus := Modulus
	hashers[0] = modulus
	seen[modulus.String()] = true

	amountOfPatterns := 3

	// patternX based hashes
	for i := 1; i < amount; i += amountOfPatterns {
		var multipier, coefficient int

		if i == 1 {
			multipier = 1
			coefficient = 0
		} else {
			multipier = toOdd(i % amount)
			coefficient = randSign(i % amount / 2)
		}

		hashFunc := NewPatternX(multipier, coefficient)
		if seen[hashFunc.String()] {
			continue
		}
		seen[hashFunc.String()] = true
		hashers[i] = hashFunc
	}

	// maxInt based hashes
	for i := 2; i < amount; i += amountOfPatterns {
		var multipier int

		k := i - 1

		if k == 1 {
			multipier = 1
		} else {
			multipier = toOdd(k % amount)
		}

		hashFunc := NewAnd(multipier)
		if seen[hashFunc.String()] {
			continue
		}
		seen[hashFunc.String()] = true
		hashers[i] = hashFunc
	}

	// bitShift based hashes
	for i := 3; i < amount; i += amountOfPatterns {
		var multipier, ander int

		k := i - 2

		if k == 1 {
			multipier = 1
		} else {
			multipier = toOdd(k % amount)
		}

		ander = k % amount / 2

		hashFunc := NewBitShift(multipier, ander)
		if seen[hashFunc.String()] {
			continue
		}
		seen[hashFunc.String()] = true
		hashers[i] = hashFunc
	}

	return hashers
}

// Modulus is a simple modulus based hash function.
var Modulus = &Hasher{
	hf: func(x, numBuckets int) int {
		return x % numBuckets
	},
	t: "x % numBuckets",
}

// NewPatternX creates new hash function with provided multipier
// and coefficient based on pattern X.
func NewPatternX(multipier, coefficient int) *Hasher {
	return &Hasher{
		hf: func(x, numBuckets int) int {
			return (multipier*x + coefficient) % numBuckets
		},
		t: fmt.Sprintf("(%d * x + %d) mod numBuckets", multipier, coefficient),
	}
}

// NewAnd creates new hash function with provided multipier which applies bitwise AND.
func NewAnd(multipier int) *Hasher {
	return &Hasher{
		hf: func(x, numBuckets int) int {
			return int(math.Abs(float64((multipier*x + x&math.MaxInt32) % numBuckets)))
		},
		t: fmt.Sprintf("(%d * x + x & maxInt) mod numBuckets", multipier),
	}
}

// NewBitShift creates new hash function with provided multipier
// and "ANDer" that utilizes bitshift under the hood.
func NewBitShift(multipier, ander int) *Hasher {
	return &Hasher{
		hf: func(x, numBuckets int) int {
			return (((x * multipier) >> 28) & ander) % numBuckets
		},
		t: fmt.Sprintf("(((x * %d) >> 28) & %d) mod numBuckets", multipier, ander),
	}
}

func toOdd(k int) int {
	return 2*k + 1
}

func randSign(k int) int {
	r := rand.Intn(2)
	if r == 0 {
		return k
	}
	return -k
}
