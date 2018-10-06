package lsh

import (
	"math"
	"sort"
)

// SetsMatrix ...
type SetsMatrix struct {
	m       map[string][]bool
	rowsNum int
	setsNum int
}

// SetsComputeMatrix optimised representation of SetsMatrix,
// optimisation is in the way it stores values - two dimensional array of booleans
// is smaller than map with strings as keys and arrays of booleans as values.
type SetsComputeMatrix struct {
	m       [][]bool
	rowsNum int
	setsNum int
}

// ToSetsMatrix returns unsorted matrix of shingles to sets.
func ToSetsMatrix(shingles [][]string) *SetsMatrix {
	m := make(map[string][]bool)

	setsNum := len(shingles)

	// iterate over provided sets of shingles
	for c, set := range shingles {
		for _, sh := range set {
			// check if matrix representation of sets has row for shingle
			// if it doesn't, then create it
			_, ok := m[sh]
			if !ok {
				m[sh] = make([]bool, setsNum)
			}
			// set 1 (true) for row with key==sh and column==c,
			// where "sh" is a shingle and "c" is a column of corresponding set
			m[sh][c] = true
		}
	}

	return &SetsMatrix{
		m:       m,
		rowsNum: len(m),
		setsNum: setsNum,
	}
}

// ToSetsComputeMatrix returns optimised sorted matrix of shingles to sets,
// where instead of a shingle itself it stores it's index.
func ToSetsComputeMatrix(shingles [][]string) *SetsComputeMatrix {
	setsMatrix := ToSetsMatrix(shingles)

	// Sort shingles for comparison consistency
	keys := make([]string, setsMatrix.rowsNum)
	i := 0
	for key := range setsMatrix.m {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	// Build optimised (only booleans) compute matrix
	m := make([][]bool, setsMatrix.rowsNum)
	for i, key := range keys {
		m[i] = make([]bool, setsMatrix.setsNum)
		for k, column := range setsMatrix.m[key] {
			m[i][k] = column
		}
	}
	return &SetsComputeMatrix{
		m:       m,
		rowsNum: setsMatrix.rowsNum,
		setsNum: setsMatrix.setsNum,
	}
}

// Minhash ...
func Minhash(shingles [][]string, numHashes int) [][]float64 {
	setsMatrix := ToSetsComputeMatrix(shingles)

	minhash := make([][]float64, numHashes)
	for i := 0; i < numHashes; i++ {
		minhash[i] = make([]float64, setsMatrix.setsNum)
		for k := 0; k < setsMatrix.setsNum; k++ {
			minhash[i][k] = math.NaN()
		}
	}

	for rNum, row := range setsMatrix.m {
		for cNum, column := range row {
			if column {
				for i := 0; i < numHashes; i++ {
					h := hashes[i](float64(rNum))
					if math.IsNaN(minhash[i][cNum]) || minhash[i][cNum] > h {
						minhash[i][cNum] = h
					}
				}
			}
		}
	}

	return minhash
}
