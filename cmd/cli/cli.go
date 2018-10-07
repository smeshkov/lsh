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

	shingleSets := make([][]string, 0)
	for _, s := range sources {
		shingleSets = append(shingleSets, getShingles(s))
	}

	signatureMatrix := lsh.Minhash(shingleSets, 3)
	fmt.Printf("signature matrix:\n%s\n\n", signatureMatrix)

	bandBuckets := lsh.LSH(signatureMatrix, 20)
	candidates := bandBuckets.FindCandidates()
	fmt.Printf("found %d candidate pair(s):\n%v\n\n", len(candidates.Index), candidates.Index)
}

func getShingles(source string) []string {
	reader, err := inout.New(source)
	if err != nil {
		fmt.Printf("can't read source %s: %v\n", source, err)
		return []string{}
	}
	lines, err := reader.ReadLines()
	if err != nil {
		fmt.Printf("can't fetch contents: %v\n", err)
	}

	textLines, _ := processor.ParseHTML(lines, false, false, false)

	return lsh.Shingle(textLines)
}
