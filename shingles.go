package lsh

import (
	"regexp"
	"strings"

	"github.com/zoomio/stopwords"
)

var (
	punctuationMarks = regexp.MustCompile(`[.,:;?!]+`)
)

type shingler struct {
	shingles   []string
	candidates [][]string
	seen       map[string]bool
}

func newShingler() *shingler {
	return &shingler{
		shingles:   make([]string, 0),
		candidates: make([][]string, 0),
		seen:       make(map[string]bool),
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
			candidate := strings.Join(sh.candidates[i], " ")
			// append result to candidates only if not seen before
			if s, ok := sh.seen[candidate]; !ok && !s {
				sh.shingles = append(sh.shingles, candidate)
				sh.seen[candidate] = true
			}
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

	seen := make(map[string]bool)

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
				_, err := sb.WriteRune(char)
				if err != nil {
					// unexpected -> panic
					panic(err)
				}

				if sb.Len() == k {
					shingles, candidates = updateShingles(i, sb.String(), seen, candidates, shingles)

					candidatesLen--
				} else {
					i++
				}
			}
		}
	}

	return shingles
}

func updateShingles(i int, candidate string, seenCandidates map[string]bool, candidates []*strings.Builder,
	shingles []string) ([]string, []*strings.Builder) {
	// append result to candidates only if it is not seen before
	if s, ok := seenCandidates[candidate]; !ok && !s {
		shingles = append(shingles, candidate)
		seenCandidates[candidate] = true
	}
	// delete from candidates
	if i == len(candidates)-1 {
		candidates = append([]*strings.Builder(nil), candidates[:i]...)
	} else {
		candidates = append(candidates[:i], candidates[i+1:]...)
	}

	return shingles, candidates
}

func isPunctuationMark(char rune) bool {
	return punctuationMarks.MatchString(string((char)))
}

func removePunctuationMarks(s string) string {
	return punctuationMarks.ReplaceAllString(s, "")
}
