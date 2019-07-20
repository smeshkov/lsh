package lsh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	aShingles = Shingle([]string{"A spokesperson for the Sudzo Corporation revealed today that studies have shown it is good for people to buy Sudzo products."})
	bShingles = Shingle([]string{"The Sudzo Corporation has revealed today that buying Sudzo products is good for people."})
	cShingles = Shingle([]string{"A spokesperson from the Sudzo Corporation has made an announcement about products of corporation."})

	simpleShingles = [][]string{
		0: {"a", "d"},
		1: {"c"},
		2: {"b", "d", "e"},
		3: {"a", "c", "d"},
	}
)

func Test_LSH_equalCandidates(t *testing.T) {
	equalShingles := [][]string{
		0: {"A spokesperson for the Sudzo Corporation revealed today that studies have shown it is good for people to buy Sudzo products."},
		1: {"There was a boy whos name was Jim. And all the friends were very good to him."},
		2: {"A spokesperson for the Sudzo Corporation revealed today that studies have shown it is good for people to buy Sudzo products."},
	}

	buckets := LSH(Minhash(equalShingles, 5), 1)
	candidates := buckets.FindCandidates()

	candidatesOf0 := candidates.GetByKey(0)
	candidatesOf2 := candidates.GetByKey(2)

	assert.Equal(t, 1, len(candidatesOf0))
	assert.Equal(t, 1, len(candidatesOf2))

	assert.Equal(t, 2, candidatesOf0[0].Index)
	assert.Equal(t, 1, candidatesOf0[0].Elections)
	assert.Equal(t, 0, candidatesOf2[0].Index)
	assert.Equal(t, 1, candidatesOf2[0].Elections)
}

func Test_LSH_equalCandidatePairs(t *testing.T) {
	equalShingles := [][]string{
		0: {"A spokesperson for the Sudzo Corporation revealed today that studies have shown it is good for people to buy Sudzo products."},
		1: {"There was a boy whos name was Jim. And all the friends were very good to him."},
		2: {"A spokesperson for the Sudzo Corporation revealed today that studies have shown it is good for people to buy Sudzo products."},
	}

	buckets := LSH(Minhash(equalShingles, 5), 1)
	candidatePairs := buckets.FindCandidatePairs()

	assert.Equal(t, 1, len(candidatePairs.Index))

	pair, ok := candidatePairs.Index["0_2"]
	assert.True(t, ok)
	assert.Equal(t, 0, pair.A)
	assert.Equal(t, 2, pair.B)
}

func Test_LSH_similarCandidates(t *testing.T) {
	similarShingles := [][]string{
		0: aShingles,
		1: {"There was a boy whos name was Jim. And all the friends were very good to him."},
		2: bShingles,
	}

	buckets := LSH(Minhash(similarShingles, 5), 3)
	candidates := buckets.FindCandidates()

	candidatesOf0 := candidates.GetByKey(0)
	candidatesOf2 := candidates.GetByKey(2)

	assert.Equal(t, 1, len(candidatesOf0))
	assert.Equal(t, 1, len(candidatesOf2))

	assert.Equal(t, 2, candidatesOf0[0].Index)
	assert.Equal(t, 1, candidatesOf0[0].Elections)
	assert.Equal(t, 0, candidatesOf2[0].Index)
	assert.Equal(t, 1, candidatesOf2[0].Elections)
}

func Test_LSH_similarCandidatePairs(t *testing.T) {
	similarShingles := [][]string{
		0: aShingles,
		1: {"There was a boy whos name was Jim. And all the friends were very good to him."},
		2: bShingles,
	}

	buckets := LSH(Minhash(similarShingles, 5), 3)
	candidatePairs := buckets.FindCandidatePairs()

	assert.Equal(t, 1, len(candidatePairs.Index))

	pair, ok := candidatePairs.Index["0_2"]
	assert.True(t, ok)
	assert.Equal(t, 0, pair.A)
	assert.Equal(t, 2, pair.B)
}
