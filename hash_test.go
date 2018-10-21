package lsh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generatesUniqueHashers(t *testing.T) {
	hashers := GenerateHashers(100)

	// make sure we generated 100 hashers
	assert.Len(t, hashers, 100)

	seen := make(map[string]bool)
	for _, value := range hashers {
		key := value.String()
		if seen[key] {
			t.Errorf("generated duplicate %s", key)
		}
		seen[key] = true
	}
}

func Test_PatternX(t *testing.T) {
	hash1 := NewPatternX(1, 1)

	assert.Equal(t, 1, hash1.Hash()(0, 5))
	assert.Equal(t, 2, hash1.Hash()(1, 5))
	assert.Equal(t, 3, hash1.Hash()(2, 5))
	assert.Equal(t, 4, hash1.Hash()(3, 5))
	assert.Equal(t, 0, hash1.Hash()(4, 5))

	hash2 := NewPatternX(3, 1)

	assert.Equal(t, 1, hash2.Hash()(0, 5))
	assert.Equal(t, 4, hash2.Hash()(1, 5))
	assert.Equal(t, 2, hash2.Hash()(2, 5))
	assert.Equal(t, 0, hash2.Hash()(3, 5))
	assert.Equal(t, 3, hash2.Hash()(4, 5))
}
