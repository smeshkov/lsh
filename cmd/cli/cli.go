package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/zoomio/inout"
	"github.com/zoomio/tagify/processor"

	"github.com/smeshkov/lsh"
)

func main() {
	source := flag.String("s", "", "List of sources separated by comma")
	numHashes := flag.Int("nh", 100, "Number of hash functions")
	verbose := flag.Bool("v", false, "Verbose")
	flag.Parse()

	if *source == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	sources := strings.Split(*source, ",")
	if len(sources) < 2 {
		fmt.Println("need at least 2 documents to find similiarities")
		os.Exit(0)
	}
	if *verbose {
		fmt.Printf("processing sources: %v\n", sources)
	}

	shingleSets := make([][]string, 0)
	for _, s := range sources {
		shingles := getShingles(s)
		// skip empty
		if len(shingles) == 0 {
			continue
		}
		shingleSets = append(shingleSets, shingles)
	}

	if len(shingleSets) < 2 {
		fmt.Printf("nothing to compare, got %d shingle set(s)\n", len(shingleSets))
		os.Exit(0)
	}

	if *verbose {
		fmt.Printf("shingle sets matrix:\n%v\n\n", shingleSets)
	}

	signatureMatrix := lsh.Minhash(shingleSets, *numHashes)
	if *verbose {
		fmt.Printf("signature matrix:\n%s\n\n", signatureMatrix)
	}

	bandBuckets := lsh.LSH(signatureMatrix, 1)
	candidates := bandBuckets.FindCandidates()

	fmt.Printf("found %d candidate pair(s)\n", len(candidates.Index))
	if len(candidates.Index) > 0 {
		fmt.Printf("%v\n\n", candidates.Index)
	}
}

func getShingles(source string) []string {
	reader, err := inout.New(source)
	if err != nil {
		fmt.Printf("can't read source %s: %v\n", source, err)
		return []string{}
	}
	lines, err := reader.ReadLines()
	if err != nil {
		fmt.Printf("can't fetch contents of %s: %v\n", source, err)
		return []string{}
	}

	textLines, _ := processor.ParseHTML(lines, false, false, false)

	return lsh.Shingle(textLines)
}
