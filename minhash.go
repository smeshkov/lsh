package lsh

import (
	"math"
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

// MinhashMatrix ...
type MinhashMatrix struct {
	m [][]int
}

// ToSetsMatrix ...
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
			// set 1 (true) for row with key==sh and column==c (column of corresponding set)
			m[sh][c] = true
		}
	}

	return &SetsMatrix{
		m:       m,
		rowsNum: len(m),
		setsNum: setsNum,
	}
}

// ToSetsComputeMatrix ...
func ToSetsComputeMatrix(shingles [][]string) *SetsComputeMatrix {
	setsMatrix := ToSetsMatrix(shingles)
	m := make([][]bool, setsMatrix.rowsNum)
	i := 0
	for _, row := range setsMatrix.m {
		m[i] = make([]bool, setsMatrix.setsNum)
		for k, column := range row {
			m[i][k] = column
		}
		i++
	}
	return &SetsComputeMatrix{
		m:       m,
		rowsNum: setsMatrix.rowsNum,
		setsNum: setsMatrix.setsNum,
	}
}

// Minhash ...
func Minhash(setsMatrix *SetsComputeMatrix) [][]float64 {

	minhash := make([][]float64, 2)
	for i := 0; i < 2; i++ {
		minhash[i] = make([]float64, setsMatrix.setsNum)
		for k := 0; k < setsMatrix.setsNum; k++ {
			minhash[i][k] = math.NaN()
		}
	}

	rNum := 0
	for _, row := range setsMatrix.m {
		for cNum, column := range row {
			if column {
				h1 := hash1(float64(rNum))
				if math.IsNaN(minhash[0][cNum]) || minhash[0][cNum] > h1 {
					minhash[0][cNum] = h1
				}

				h2 := hash2(float64(rNum))
				if math.IsNaN(minhash[1][cNum]) || minhash[1][cNum] > h2 {
					minhash[1][cNum] = h2
				}
			}
		}
		rNum++
	}

	return minhash
}

func hash1(x float64) float64 {
	return math.Mod(x+1, 5)
}

func hash2(x float64) float64 {
	return math.Mod(3*x+1, 5)
}
