package lsh

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// SetsMatrix contains index of shingles to sets,
// i.e. each key in the map is a string representation of shingle
// and a value is a list of booleans corresponding to the documents,
// true value in the list means that document contains the shingle.
type SetsMatrix struct {
	m       map[string][]bool
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
			// where "sh" is a shingle and "c" is a column/index of corresponding set/document
			m[sh][c] = true
		}
	}

	return &SetsMatrix{
		m:       m,
		setsNum: setsNum,
	}
}

// ShinglesNum returns number of shingles.
func (sm *SetsMatrix) ShinglesNum() int {
	return len(sm.m)
}

// Clone makes a copy of this SetsMatrix.
func (sm *SetsMatrix) Clone() *SetsMatrix {
	m := make(map[string][]bool)
	for k, v := range sm.m {
		m[k] = v
	}
	return &SetsMatrix{
		m:       m,
		setsNum: sm.setsNum,
	}
}

// SetsComputeMatrix optimised representation of SetsMatrix,
// optimisation is in the way it stores values - two dimensional array of booleans
// is smaller than map with strings as keys and arrays of booleans as values.
//
// How optimisation is built?
// Shingles are being sorted, so each shingle has certain ordered position,
// therefore there is no more need in storing actual value of the shingle,
// because it's position is just enough.
type SetsComputeMatrix struct {
	m       [][]bool
	rowsNum int
	setsNum int
}

// ToSetsComputeMatrix returns optimised sorted matrix of sets matrix,
// where instead of a shingle itself it stores it's index.
func ToSetsComputeMatrix(setsMatrix *SetsMatrix) *SetsComputeMatrix {
	rowsNum := setsMatrix.ShinglesNum()
	// Sort shingles for comparison consistency
	keys := make([]string, rowsNum)
	i := 0
	for key := range setsMatrix.m {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	// Build optimised (only booleans) compute matrix
	m := make([][]bool, rowsNum)
	for i, key := range keys {
		m[i] = make([]bool, setsMatrix.setsNum)
		copy(m[i], setsMatrix.m[key])
	}
	return &SetsComputeMatrix{
		m:       m,
		rowsNum: rowsNum,
		setsNum: setsMatrix.setsNum,
	}
}

func (scm *SetsComputeMatrix) String() string {
	var sb strings.Builder
	for i, row := range scm.m {
		for j, column := range row {
			if j > 0 {
				checkWriteStringError(sb.WriteString(","))
			}
			if column {
				checkWriteStringError(sb.WriteString("1"))
			} else {
				checkWriteStringError(sb.WriteString("0"))
			}
		}
		if i < len(scm.m)-1 {
			checkWriteStringError(sb.WriteString("\n"))
		}
	}
	return sb.String()
}

// SignatureMatrix - each row represents a hash function and each column represents a set.
type SignatureMatrix [][]float64

func (sm SignatureMatrix) String() string {
	var sb strings.Builder
	for i, row := range sm {
		for j, column := range row {
			if j > 0 {
				checkWriteStringError(sb.WriteString(","))
			}
			checkWriteStringError(sb.WriteString(fmt.Sprintf("%.0f", column)))
		}
		if i < len(sm)-1 {
			checkWriteStringError(sb.WriteString("\n"))
		}
	}
	return sb.String()
}

// Minhash performs minhashing operations on the given shingles,
// with the given number (`numHashes`) of generated hashes functions.
func Minhash(shingles [][]string, numHashes int) SignatureMatrix {
	return MinhashWithHashers(shingles, GenerateHashers(numHashes))
}

// MinhashWithHashers performs minhashing operations on the given shingles,
// with the given hashes functions.
func MinhashWithHashers(shingles [][]string, hashers []*Hasher) SignatureMatrix {
	return minhashSetsMatrix(ToSetsMatrix(shingles), hashers)
}

func minhashSetsMatrix(setsMatrix *SetsMatrix, hashers []*Hasher) SignatureMatrix {
	setsComputeMatrix := ToSetsComputeMatrix(setsMatrix)
	numHashes := len(hashers)

	// build a signature matrix, initialy by filing all the values with NaN.
	minhash := make(SignatureMatrix, numHashes)
	for i := 0; i < numHashes; i++ {
		minhash[i] = make([]float64, setsComputeMatrix.setsNum)
		for k := 0; k < setsComputeMatrix.setsNum; k++ {
			minhash[i][k] = math.NaN()
		}
	}

	// run through compute matrix and perform hashing,
	// if set (document) has shingle represented in it.
	for rNum, row := range setsComputeMatrix.m {
		for cNum, column := range row {
			if column {
				for i := 0; i < numHashes; i++ {
					h := hashers[i].Hash()(rNum, setsComputeMatrix.rowsNum)
					if math.IsNaN(minhash[i][cNum]) || minhash[i][cNum] > float64(h) {
						minhash[i][cNum] = float64(h)
					}
				}
			}
		}
	}

	return minhash
}

func checkWriteStringError(ignored int, err error) {
	if err != nil {
		panic(fmt.Sprintf("error in building a string from SetsComputeMatrix: %v", err))
	}
}
