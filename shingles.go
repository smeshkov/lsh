package lsh

import (
	"regexp"
	"strings"

	"github.com/zoomio/stopwords"
)

const (
	defaultK = 5
)

var (
	punctuationMarks = regexp.MustCompile(`[.,:;?!]+`)
)

type shingler struct {
	shingles   []string
	candidates [][]string
}

func newShingler() *shingler {
	return &shingler{
		shingles:   make([]string, 0),
		candidates: make([][]string, 0),
	}
}

func (sh *shingler) appendCandidate() {
	sh.candidates = append(sh.candidates, make([]string, 0))
}

func (sh *shingler) appendWord(word string) {
	candidatesLen := len(sh.candidates)

	if candidatesLen == 0 {
		return
	}

	i := 0
	for {
		if i == candidatesLen {
			break
		}

		sh.candidates[i] = append(sh.candidates[i], word)

		if len(sh.candidates[i]) == 3 {
			// append to result shingles
			sh.shingles = append(sh.shingles, strings.Join(sh.candidates[i], " "))
			// delete from candidates
			if i == candidatesLen-1 {
				sh.candidates = append([][]string(nil), sh.candidates[:i]...)
			} else {
				sh.candidates = append(sh.candidates[:i], sh.candidates[i+1:]...)
			}

			candidatesLen--
		} else {
			i++
		}
	}
}

// Shingle produces shingles of a stop word followed by
// the next two words from the given lines of strings.
func Shingle(lines []string) []string {
	sh := newShingler()

	for _, line := range lines {

		words := strings.Fields(line)
		for _, word := range words {
			w := removePunctuationMarks(word)
			if stopwords.IsStopWord(strings.ToLower(w)) {
				sh.appendCandidate()
			}
			sh.appendWord(w)
		}
	}

	return sh.shingles
}

// KShingle produces shingles of given size k.
func KShingle(lines []string, k int) []string {
	shingles := make([]string, 0)
	candidates := make([]*strings.Builder, 0)

	for _, line := range lines {
		for _, char := range line {
			if isPunctuationMark(char) {
				continue
			}
			candidates = append(candidates, &strings.Builder{})
			candidatesLen := len(candidates)

			i := 0
			for {
				if i == candidatesLen {
					break
				}

				sb := candidates[i]
				sb.WriteRune(char)

				if sb.Len() == k {
					shingles = append(shingles, sb.String())
					// delete from candidates
					if i == candidatesLen-1 {
						candidates = append([]*strings.Builder(nil), candidates[:i]...)
					} else {
						candidates = append(candidates[:i], candidates[i+1:]...)
					}

					candidatesLen--
				} else {
					i++
				}
			}
		}
	}

	return shingles
}

func isPunctuationMark(char rune) bool {
	return punctuationMarks.MatchString(string((char)))
}

func removePunctuationMarks(s string) string {
	return punctuationMarks.ReplaceAllString(s, "")
}
