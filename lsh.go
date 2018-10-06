package lsh

type address struct {
	bandNum int
	setNum  int
}

type candidateBuckets [][]*address

// BandBuckets ...
type BandBuckets struct {
	bandBuckets []candidateBuckets
	hFunc       Hash
}

func newBandedBuckets(bands, buckets int, hFunc Hash) *BandBuckets {
	bb := &BandBuckets{
		bandBuckets: make([]candidateBuckets, bands),
		hFunc:       hFunc,
	}
	for index := 0; index < bands; index++ {
		bb.bandBuckets[index] = make(candidateBuckets, buckets)
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
	buckets := bb.bandBuckets[bandNum]

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
	bb.bandBuckets[bandNum] = buckets
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
