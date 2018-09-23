package lsh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	aShingles = Shingle([]string{"A spokesperson for the Sudzo Corporation revealed today that studies have shown it is good for people to buy Sudzo products."})
	bShingles = Shingle([]string{"The Sudzo Corporation has revealed today that buying Sudzo products is good for people."})
	cShingles = Shingle([]string{"A spokesperson from the Sudzo Corporation has made an announcement about products of corporation."})
)

func Test_ToSetsMatrix(t *testing.T) {
	shingles := make([][]string, 3)
	shingles[0] = aShingles
	shingles[1] = bShingles
	shingles[2] = cShingles

	setsMatrix := ToSetsMatrix(shingles)

	intersection := make(map[string]bool)
	for _, set := range shingles {
		for _, sh := range set {
			intersection[sh] = true
		}
	}

	// Assert that number of sets in setsMatrix is the amount of provided sets of shingles
	assert.Equal(t, setsMatrix.setsNum, len(shingles))

	// Assert that number of rows in setsMatrix is the length of the intersection of provided sets of shingles
	assert.Equal(t, setsMatrix.rowsNum, len(intersection))

	// Assert that 1st and 2nd sets have the same entry for shingle - "is good for"
	assert.True(t, setsMatrix.m["is good for"][0])
	assert.True(t, setsMatrix.m["is good for"][1])

	// Assert that 1st and 3rd sets have the same entry for shingle - "the Sudzo Corporation"
	assert.True(t, setsMatrix.m["the Sudzo Corporation"][0])
	assert.True(t, setsMatrix.m["the Sudzo Corporation"][2])
}

func Test_ToSetsMatrix2(t *testing.T) {
	shingles := make([][]string, 4)
	shingles[0] = []string{"a", "d"}
	shingles[1] = []string{"c"}
	shingles[2] = []string{"b", "d", "e"}
	shingles[3] = []string{"a", "c", "d"}

	setsMatrix := ToSetsMatrix(shingles)

	// Assert that number of sets in setsMatrix is the amount of provided sets of shingles
	assert.Equal(t, setsMatrix.setsNum, 4)

	// Assert that number of rows in setsMatrix is the length of the intersection of provided sets of shingles
	assert.Equal(t, setsMatrix.rowsNum, 5)

	// "a" is in set1 and set4
	assert.True(t, setsMatrix.m["a"][0])
	assert.True(t, setsMatrix.m["a"][3])

	// "b" is only in set3
	assert.True(t, setsMatrix.m["b"][2])

	// "c" is in set2 and set4
	assert.True(t, setsMatrix.m["c"][1])
	assert.True(t, setsMatrix.m["c"][3])

	// "d" is in set1, set3 and set4
	assert.True(t, setsMatrix.m["d"][0])
	assert.True(t, setsMatrix.m["d"][2])
	assert.True(t, setsMatrix.m["d"][3])

	// "e" is only in set3
	assert.True(t, setsMatrix.m["e"][2])
}

func Test_MinHash(t *testing.T) {
	setsMatrix := &SetsComputeMatrix{
		m: [][]bool{
			0: []bool{true, false, false, true},
			1: []bool{false, false, true, false},
			2: []bool{false, true, false, true},
			3: []bool{true, false, true, true},
			4: []bool{false, false, true, false},
		},
		rowsNum: 5,
		setsNum: 4,
	}

	minhash := Minhash(setsMatrix)

	// h1 hasing function assertions
	assert.Equal(t, 1.0, minhash[0][0])
	assert.Equal(t, 3.0, minhash[0][1])
	assert.Equal(t, 0.0, minhash[0][2])
	assert.Equal(t, 1.0, minhash[0][3])

	// h2 hasing function assertions
	assert.Equal(t, 0.0, minhash[1][0])
	assert.Equal(t, 2.0, minhash[1][1])
	assert.Equal(t, 0.0, minhash[1][2])
	assert.Equal(t, 0.0, minhash[1][3])
}

func Test_hash1(t *testing.T) {
	assert.Equal(t, 1.0, hash1(0))
	assert.Equal(t, 2.0, hash1(1))
	assert.Equal(t, 3.0, hash1(2))
	assert.Equal(t, 4.0, hash1(3))
	assert.Equal(t, 0.0, hash1(4))
}

func Test_hash2(t *testing.T) {
	assert.Equal(t, 1.0, hash2(0))
	assert.Equal(t, 4.0, hash2(1))
	assert.Equal(t, 2.0, hash2(2))
	assert.Equal(t, 0.0, hash2(3))
	assert.Equal(t, 3.0, hash2(4))
}
