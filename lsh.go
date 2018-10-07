package lsh

import (
	"fmt"
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

type address struct {
	bandNum int
	setNum  int
}

type candidateBuckets [][]*address

// BandBuckets stores candidates for comparisson in the same bucket,
// bucket groups are separated by band.
type BandBuckets struct {
	bands []candidateBuckets
	hFunc Hash
}

func newBandedBuckets(bands, buckets int, hFunc Hash) *BandBuckets {
	bb := &BandBuckets{
		bands: make([]candidateBuckets, bands),
		hFunc: hFunc,
	}
	for index := 0; index < bands; index++ {
		bb.bands[index] = make(candidateBuckets, buckets)
	}
	return bb
}

func (bb *BandBuckets) hashToBucket(vector []float64, bandNum, setNum int) {
	// transform vector to int
	h := 0.0
	for _, v := range vector {
		h = h*10 + v
	}

	// hash transformed int
	h = bb.hFunc(h)

	// get buckets for given band number
	buckets := bb.bands[bandNum]

	// get backet addres for candidate
	bucketNum := int(h) % len(buckets)

	// get backet for candidate
	cb := buckets[bucketNum]

	// append candidate details to bucket
	cb = append(cb, &address{
		bandNum: bandNum,
		setNum:  setNum,
	})

	// put apdated list of candidates into bucket
	buckets[bucketNum] = cb

	// update bands with updated buckets
	bb.bands[bandNum] = buckets
}

// FindCandidates provides slice of candidate groups,
// i.e. each entry in the slice is the list of candidates that ended up in the same bucket.
func (bb *BandBuckets) FindCandidates() *Candidates {
	candidates := &Candidates{Index: make(map[string]*CandidatePair)}

	// iterate through bands and pick up candidates for comparison from each of the bands
	for _, band := range bb.bands {
		for _, bucket := range band {
			// skip buckets with none or single candidate
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
	numBuckets := numSets
	numRows := numHashes / bands

	bb := newBandedBuckets(bands, numBuckets, hashes[3])

	for b := 0; b < bands; b++ {
		bandVectors := make([][]float64, numSets)
		for h := b * numRows; h < numHashes; h++ {
			for s := 0; s < numSets; s++ {
				bandVectors[s] = append(bandVectors[s], signatureMatrix[h][s])
			}
		}
		for i, vector := range bandVectors {
			bb.hashToBucket(vector, b, i)
		}
	}

	return bb
}
