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
	assert.Equal(t, 1.0, hash1(0))
	assert.Equal(t, 2.0, hash1(1))
	assert.Equal(t, 3.0, hash1(2))
	assert.Equal(t, 4.0, hash1(3))
	assert.Equal(t, 0.0, hash1(4))
}

func Test_hash2(t *testing.T) {
	assert.Equal(t, 1.0, hash2(0))
	assert.Equal(t, 4.0, hash2(1))
	assert.Equal(t, 2.0, hash2(2))
	assert.Equal(t, 0.0, hash2(3))
	assert.Equal(t, 3.0, hash2(4))
}
