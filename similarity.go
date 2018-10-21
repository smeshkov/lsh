package lsh

// Jaccard ...
func Jaccard(a, b []string) float64 {
	setA := make(map[string]bool)
	setB := make(map[string]bool)
	union := make(map[string]bool)

	// fill sets
	for i := 0; i < len(a) || i < len(b); i++ {
		if i < len(a) {
			setA[a[i]] = true
			union[a[i]] = true
		}
		if i < len(b) {
			setB[b[i]] = true
			union[b[i]] = true
		}
	}

	return float64(len(intersetion(setA, setB))) / float64(len(union))
}

func intersetion(setA, setB map[string]bool) map[string]bool {
	intersection := make(map[string]bool)
	// loop over smaller set
	if len(setA) < len(setB) {
		for key := range setA {
			if setB[key] {
				intersection[key] = true
			}
		}
	} else {
		for key := range setB {
			if setA[key] {
				intersection[key] = true
			}
		}
	}
	return intersection
}
