package lsh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	aText = "A spokesperson for the Sudzo Corporation revealed today that studies have shown it is good for people to buy Sudzo products."
	dupedText = "A spokesperson for the Sudzo Corporation. A spokesperson for the Sudzo Corporation."
)

func Test_isPunctuationMark(t *testing.T) {
	assert.True(t, isPunctuationMark('.'))
	assert.True(t, isPunctuationMark(','))
	assert.True(t, isPunctuationMark(':'))
	assert.True(t, isPunctuationMark(';'))
	assert.True(t, isPunctuationMark('?'))
	assert.True(t, isPunctuationMark('!'))
	assert.False(t, isPunctuationMark(' '))
	assert.False(t, isPunctuationMark('a'))
	assert.False(t, isPunctuationMark('1'))
	assert.False(t, isPunctuationMark('#'))
}

func Test_removePunctuationMark(t *testing.T) {
	assert.Equal(t, "for people to buy Sudzo products", removePunctuationMarks("for people to buy Sudzo products."))
	assert.Equal(t, " for people to buy Sudzo products", removePunctuationMarks("... for people to buy Sudzo products."))
	assert.Equal(t, "for people to buy Sudzo products", removePunctuationMarks("...for people to buy Sudzo products."))
	assert.Equal(t, "Hello world", removePunctuationMarks("Hello, world!"))
}

func Test_Shingle(t *testing.T) {
	shingles := Shingle([]string{aText})

	assert.Len(t, shingles, 10)
	assert.Equal(t, "A spokesperson for", shingles[0])
	assert.Equal(t, "for the Sudzo", shingles[1])
	assert.Equal(t, "the Sudzo Corporation", shingles[2])
	assert.Equal(t, "that studies have", shingles[3])
	assert.Equal(t, "have shown it", shingles[4])
	assert.Equal(t, "shown it is", shingles[5])
	assert.Equal(t, "it is good", shingles[6])
	assert.Equal(t, "is good for", shingles[7])
	assert.Equal(t, "for people to", shingles[8])
	assert.Equal(t, "to buy Sudzo", shingles[9])
}

func Test_Shingle_DeDuplicates(t *testing.T) {
	shingles := Shingle([]string{dupedText})

	assert.Len(t, shingles, 3)
	assert.Equal(t, "A spokesperson for", shingles[0])
	assert.Equal(t, "for the Sudzo", shingles[1])
	assert.Equal(t, "the Sudzo Corporation", shingles[2])
}

func Test_KShingle(t *testing.T) {
	shingles := KShingle([]string{aText}, 9)

	assert.Len(t, shingles, 115)
	assert.Equal(t, "A spokesp", shingles[0])
	assert.Equal(t, " spokespe", shingles[1])
	assert.Equal(t, "spokesper", shingles[2])
	// ...
	assert.Equal(t, "zo produc", shingles[len(shingles)-3])
	assert.Equal(t, "o product", shingles[len(shingles)-2])
	assert.Equal(t, " products", shingles[len(shingles)-1])
}

func Test_KShingle_DeDuplicates(t *testing.T) {
	shingles := KShingle([]string{dupedText}, 9)

	assert.Len(t, shingles, 41)
	assert.Equal(t, "A spokesp", shingles[0])
	assert.Equal(t, " spokespe", shingles[1])
	assert.Equal(t, "spokesper", shingles[2])
	// ...
	assert.Equal(t, "on A spok", shingles[len(shingles)-3])
	assert.Equal(t, "n A spoke", shingles[len(shingles)-2])
	assert.Equal(t, " A spokes", shingles[len(shingles)-1])
}
