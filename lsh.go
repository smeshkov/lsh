package lsh

import (
	"fmt"
	"math"
	"sort"
)

// CandidatePair ...
type CandidatePair struct {
	A         int    // index of a candidate A
	B         int    // index of a candidate B
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

// CandidatePairs ...
type CandidatePairs struct {
	Index map[string]*CandidatePair
}

// Put ...
func (c *CandidatePairs) Put(a, b int) {
	cp := newCandidatePair(a, b)
	if _, ok := c.Index[cp.signature]; ok {
		c.Index[cp.signature].Elections++
	} else {
		c.Index[cp.signature] = cp
	}
}

// Keys returns just keys of candidate pairs, useful for debugging or just printing to STDIN.
func (c *CandidatePairs) Keys() []string {
	var keys []string
	for key := range c.Index {
		keys = append(keys, key)
	}
	return keys
}

// Candidate ...
type Candidate struct {
	Index     int // index of candidate
	Elections int // how many times candidate ended up in the same bucket
}

// ------------------------------------------ Compare ------------------------------------------

// Cmp is the type of a "compare" function that defines the ordering of its Candidate arguments.
type Cmp func(i1, i2 *Candidate) int

// ------------------------------------------ Sort ------------------------------------------

// By is the type of a "less" function that defines the ordering of its Candidate arguments.
type By func(i1, i2 *Candidate) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(items []*Candidate) {
	sorter := &itemSorter{
		items: items,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(sorter)
}

// itemSorter joins a By function and a slice of Items to be sorted.
type itemSorter struct {
	items []*Candidate
	by    func(p1, p2 *Candidate) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *itemSorter) Len() int {
	return len(s.items)
}

// Swap is part of sort.Interface.
func (s *itemSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *itemSorter) Less(i, j int) bool {
	return s.by(s.items[i], s.items[j])
}

// Candidates is an index of candidates, keyed by document signature (e.g. document number),
// with value representing list of documents which ended up in the same bucket.
type Candidates struct {
	Index map[int]map[int]*Candidate
}

// Put puts candidate "b" to the adjacent map of candidates of "a".
func (c *Candidates) Put(a, b int) {
	_, ok := c.Index[a]
	if !ok {
		c.Index[a] = make(map[int]*Candidate)
	}
	_, ok = c.Index[a][b]
	if !ok {
		c.Index[a][b] = &Candidate{Index: b, Elections: 1}
	} else {
		c.Index[a][b].Elections++
	}
}

// GetByKey returns list of adjacent candidates for given candidate "key".
func (c *Candidates) GetByKey(key int) []*Candidate {
	if cndMap, ok := c.Index[key]; ok {
		res := make([]*Candidate, len(cndMap))
		var i int
		for _, cnd := range cndMap {
			res[i] = cnd
			i++
		}
		return res
	}
	return []*Candidate{}
}

// GetByKeySorted returns list of adjacent candidates for given candidate "key"
// sorted by Elections in descending order.
func (c *Candidates) GetByKeySorted(key int) []*Candidate {
	candidates := c.GetByKey(key)
	By(func(i1, i2 *Candidate) bool {
		return i1.Elections > i2.Elections
	}).Sort(candidates)
	return candidates
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
	candidates := &Candidates{Index: make(map[int]map[int]*Candidate)}

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
					candidates.Put(bucket[j].setNum, bucket[i].setNum)
				}
			}
		}
	}

	return candidates
}

// FindCandidatePairs provides slice of candidate groups,
// i.e. each entry in the slice is the list of candidates that ended up in the same bucket.
func (bb *BandBuckets) FindCandidatePairs() *CandidatePairs {
	candidates := &CandidatePairs{Index: make(map[string]*CandidatePair)}

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
					candidates.Put(bucket[j].setNum, bucket[i].setNum)
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
