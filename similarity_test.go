package lsh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Jaccard(t *testing.T) {
	assert.Equal(t, 1.0, Jaccard([]string{"a", "b"}, []string{"a", "b"}))
	assert.Equal(t, 0.5, Jaccard([]string{"a", "b"}, []string{"a"}))
}
