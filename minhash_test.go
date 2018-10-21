package lsh

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	// Assert that number of rows in setsMatrix is the length of the union of provided unique sets of shingles
	assert.Equal(t, setsMatrix.rowsNum, len(intersection))

	// Assert that 1st and 2nd sets have the same entry for shingle - "is good for"
	assert.True(t, setsMatrix.m["is good for"][0])
	assert.True(t, setsMatrix.m["is good for"][1])

	// Assert that 1st and 3rd sets have the same entry for shingle - "the Sudzo Corporation"
	assert.True(t, setsMatrix.m["the Sudzo Corporation"][0])
	assert.True(t, setsMatrix.m["the Sudzo Corporation"][2])
}

func Test_ToSetsMatrix2(t *testing.T) {
	setsMatrix := ToSetsMatrix(simpleShingles)

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

func Test_ToSetsComputeMatrix(t *testing.T) {

	// Input matrix:
	// row | s1 | s2 | s3 | s4
	//  0  |  1 |  0 |  0 |  1
	//  1  |  0 |  0 |  1 |  0
	//  2  |  0 |  1 |  0 |  1
	//  3  |  1 |  0 |  1 |  1
	//  4  |  0 |  0 |  1 |  0

	expected := &SetsComputeMatrix{
		m: [][]bool{
			0: {true, false, false, true},
			1: {false, false, true, false},
			2: {false, true, false, true},
			3: {true, false, true, true},
			4: {false, false, true, false},
		},
		rowsNum: 5,
		setsNum: 4,
	}

	actual := ToSetsComputeMatrix(simpleShingles)

	assert.Equal(t, expected.rowsNum, actual.rowsNum)
	assert.Equal(t, expected.setsNum, actual.setsNum)
	assert.Equal(t, len(expected.m), len(actual.m))

	for i, row := range expected.m {
		assert.Equal(t, len(row), len(actual.m[i]))
		for k, value := range row {
			assert.Equal(t, value, actual.m[i][k])
		}
	}
}

func Test_MinHash_EnforcesOrder(t *testing.T) {

	// Input matrix:
	// row | s1 | s2 | s3 | s4
	//  0  |  1 |  0 |  0 |  1
	//  1  |  0 |  0 |  1 |  0
	//  2  |  0 |  1 |  0 |  1
	//  3  |  1 |  0 |  1 |  1
	//  4  |  0 |  0 |  1 |  0

	minhash := MinhashWithHashers(simpleShingles, []*Hasher{NewPatternX(1, 1), NewPatternX(3, 1)})

	// Output matrix:
	//  h  | s1 | s2 | s3 | s4
	// h1  |  1 |  3 |  0 |  1
	// h2  |  0 |  2 |  0 |  0

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
