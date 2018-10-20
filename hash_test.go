package lsh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	hash1 = hashes[0]
	hash2 = hashes[1]
)

func Test_hash1(t *testing.T) {
	assert.Equal(t, 1, hash1(0, 5))
	assert.Equal(t, 2, hash1(1, 5))
	assert.Equal(t, 3, hash1(2, 5))
	assert.Equal(t, 4, hash1(3, 5))
	assert.Equal(t, 0, hash1(4, 5))
}

func Test_hash2(t *testing.T) {
	assert.Equal(t, 1, hash2(0, 5))
	assert.Equal(t, 4, hash2(1, 5))
	assert.Equal(t, 2, hash2(2, 5))
	assert.Equal(t, 0, hash2(3, 5))
	assert.Equal(t, 3, hash2(4, 5))
}
