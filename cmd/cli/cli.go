package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/zoomio/inout"

	"github.com/smeshkov/lsh"
)

func main() {
	source := flag.String("s", "", "Source")
	flag.Parse()

	if *source == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	sources := strings.Fields(*source)
	if len(sources) == 0 {
		fmt.Printf("wrong source: %s", *source)
		flag.PrintDefaults()
		os.Exit(0)
	}

	shingleSets := make([][]string, 0)
	for _, s := range sources {
		shingleSets = append(shingleSets, getShingles(s))
	}

	signatureMatrix := lsh.Minhash(shingleSets, 3)
	bandBuckets := lsh.LSH(signatureMatrix, 20)

	fmt.Printf("%v\n", bandBuckets)
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
	return lsh.Shingle(lines)
}
