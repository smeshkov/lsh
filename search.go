package lsh

// Search configuration options.
var (
	// Hashers sets hashers funcs.
	Hashers = func(hashers []*Hasher) SearchOption {
		return func(s *Search) {
			s.hashers = hashers
		}
	}

	// HashersNum generates desired number of Hashers for Search.
	HashersNum = func(hashersNum int) SearchOption {
		return func(s *Search) {
			s.hashers = GenerateHashers(hashersNum)
		}
	}

	// BandsNum sets number of bands.
	BandsNum = func(bandsNum int) SearchOption {
		return func(s *Search) {
			s.bandsNum = bandsNum
		}
	}

	// Index sets index for search.
	Index = func(index *SetsMatrix) SearchOption {
		return func(s *Search) {
			s.index = index
		}
	}
)

// Search ...
type Search struct {
	hashers  []*Hasher
	bandsNum int
	index    *SetsMatrix
}

// NewSearch creates new instance of Search.
func NewSearch(options ...SearchOption) Search {
	s := &Search{}

	// apply custom configuration
	for _, option := range options {
		option(s)
	}

	// set defaults if needed
	if s.hashers == nil || len(s.hashers) == 0 {
		HashersNum(100)(s)
	}
	if s.bandsNum == 0 {
		BandsNum(20)(s)
	}
	if s.index == nil {
		Index(ToSetsMatrix([][]string{}))(s)
	}

	return *s
}

// Find finds candidates for given query string.
func (s *Search) Find(query string) *Candidates {
	shingles := Shingle([]string{query})
	index := s.reIndex(shingles)
	signatureMatrix := minhashSetsMatrix(index, s.hashers)
	bandBuckets := LSH(signatureMatrix, s.bandsNum)
	return bandBuckets.FindCandidates()
	/* found := candidates.GetByKey(index.setsNum - 1)
	result := make([]string, len(found))
	for i, v := range found {
		result[i] = s.documents[v.Index]
	}
	return result */
}

// reIndex re-indexes current set of documents given the additional query.
func (s *Search) reIndex(shingles []string) *SetsMatrix {
	clone := s.index.Clone()
	setIndex := clone.setsNum
	clone.setsNum++
	for _, sh := range shingles {
		// check if matrix representation of sets has row for shingle
		// if it doesn't, then create it
		_, ok := clone.m[sh]
		if !ok {
			clone.m[sh] = make([]bool, clone.setsNum)
		}
		// set 1 (true) for row with key==sh and column==c,
		// where "sh" is a shingle and "c" is a column/index of corresponding set/document
		clone.m[sh][setIndex] = true
	}
	return clone
}

// SearchOption allows to customise configuration.
type SearchOption func(*Search)
