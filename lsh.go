package lsh

import (
	"fmt"
	"math"
)

// CandidatePair ...
type CandidatePair struct {
	A         int    // index of candidate A
	B         int    // index of candidate B
	Elections int    // how many times candidates ended up in the same bucket
	signature string // unique signature that identifies candidates
}

func newCandidatePair(a, b int) *CandidatePair {
	if a > b {
		return &CandidatePair{
			A:         b,
			B:         a,
			Elections: 1,
			signature: fmt.Sprintf("%d_%d", b, a),
		}
	}
	return &CandidatePair{
		A:         a,
		B:         b,
		Elections: 1,
		signature: fmt.Sprintf("%d_%d", a, b),
	}
}

// Candidates ...
type Candidates struct {
	Index map[string]*CandidatePair
}

// Put ...
func (c *Candidates) Put(a, b int) {
	cp := newCandidatePair(a, b)
	if _, ok := c.Index[cp.signature]; ok {
		c.Index[cp.signature].Elections++
	} else {
		c.Index[cp.signature] = cp
	}
}

// Keys returns just keys of candidate pairs, useful for debugging or just printing to STDIN.
func (c *Candidates) Keys() []string {
	var keys []string
	for key := range c.Index {
		keys = append(keys, key)
	}
	return keys
}

type address struct {
	bandNum int
	setNum  int
}

type candidateBuckets [][]*address

// BandBuckets stores candidates for comparisson in the same bucket,
// bucket groups are separated by band.
type BandBuckets struct {
	bands []candidateBuckets
}

func newBandBuckets(bands, buckets int) *BandBuckets {
	bb := &BandBuckets{
		bands: make([]candidateBuckets, bands),
	}
	for index := 0; index < bands; index++ {
		bb.bands[index] = make(candidateBuckets, buckets)
	}
	return bb
}

// hashToBucket hashes given vector into bucket and returns hash and bucket number.
func (bb *BandBuckets) hashToBucket(vector []float64, bandNum, setNum int) (int, int) {
	// hash vector to int
	var h int
	for _, v := range vector {
		h = 31*h + (int(v) & math.MaxInt32)
	}

	// get buckets for given band number
	buckets := bb.bands[bandNum]

	// hash transformed int
	// get backet addres for candidate
	bucketNum := int(math.Abs(float64(h % len(buckets))))

	// get backet for candidate
	candidateBucket := buckets[bucketNum]

	// append candidate details to bucket
	candidateBucket = append(candidateBucket, &address{
		bandNum: bandNum,
		setNum:  setNum,
	})

	// put apdated list of candidates into bucket
	buckets[bucketNum] = candidateBucket

	// update bands with updated buckets
	bb.bands[bandNum] = buckets

	return h, bucketNum
}

// FindCandidates provides slice of candidate groups,
// i.e. each entry in the slice is the list of candidates that ended up in the same bucket.
func (bb *BandBuckets) FindCandidates() *Candidates {
	candidates := &Candidates{Index: make(map[string]*CandidatePair)}

	// iterate through bands and pick up candidates for comparison from each of the bands
	for _, band := range bb.bands {
		for _, bucket := range band {
			// skip buckets with zero or one candidate
			if len(bucket) < 2 {
				continue
			}
			for i := 0; i < len(bucket); i++ {
				for j := i + 1; j < len(bucket); j++ {
					candidates.Put(bucket[i].setNum, bucket[j].setNum)
				}
			}
		}
	}

	return candidates
}

// LSH applies Locality Sesnitive Hashing (banded approach) onto the given signature matrix
// in order to find candidate pairs for similiarity.
func LSH(signatureMatrix [][]float64, bands int) *BandBuckets {
	numHashes := len(signatureMatrix)
	numSets := len(signatureMatrix[0])
	numBuckets := numHashes
	numRows := numHashes / bands

	// debug logging
	// fmt.Printf("numBands %d, numHashes %d, numSets %d, numBuckets %d, numRows in band %d\n",
	// bands, numHashes, numSets, numBuckets, numRows)

	bb := newBandBuckets(bands, numBuckets)

	for b := 0; b < bands; b++ {
		bandVectors := make([][]float64, numSets)
		bandOffset := b * numRows
		bandEnd := (b + 1) * numRows

		// debug logging
		// fmt.Printf("bandOffset %d, bandEnd %d\n",
		// 	bandOffset, bandEnd)

		for h := bandOffset; bandEnd <= numHashes && h < bandEnd; h++ {
			for s := 0; s < numSets; s++ {
				bandVectors[s] = append(bandVectors[s], signatureMatrix[h][s])
			}
		}

		// debug logging
		// fmt.Printf("bandVectors:\n%v\n\n", bandVectors)

		for i, vector := range bandVectors {
			bb.hashToBucket(vector, b, i)
		}
	}

	return bb
}
