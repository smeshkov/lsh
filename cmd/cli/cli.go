package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gpestana/htmlizer"
	"github.com/zoomio/inout"

	"github.com/smeshkov/lsh"
)

var (
	// shingle command
	shingleCmd    = flag.NewFlagSet("shingle", flag.ExitOnError)
	shingleSource = shingleCmd.String("s", "", "Sources of the text.")
	shingleK      = shingleCmd.Bool("k", false, "Enables K-shingling.")

	// LSH command
	lshCmd       = flag.NewFlagSet("lsh", flag.ExitOnError)
	lshSources   = lshCmd.String("s", "", "List of sources separated by comma.")
	lshNumHashes = lshCmd.Int("hashes", 0, "Number of hash functions.")
	lshNumBands  = lshCmd.Int("bands", 0, "Number of bands.")

	// similarity command
	simCmd       = flag.NewFlagSet("sim", flag.ExitOnError)
	simSources   = simCmd.String("s", "", "List of sources separated by comma.")
	simKShingles = simCmd.Int("k", 9, "Number of shingles for K-shingling approach.")
)

func main() {
	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	for strings.HasPrefix(cmd, "-") {
		cmd = strings.TrimPrefix(cmd, "-")
	}
	if cmd == "help" || cmd == "h" {
		printUsage()
		os.Exit(1)
	}

	switch cmd {
	case shingleCmd.Name():
		doShingles(shingleCmd)
	case lshCmd.Name():
		doLSH(lshCmd)
	case simCmd.Name():
		doSim(simCmd)
	default:
		fmt.Printf("unknown: %s\n", cmd)
		os.Exit(2)
	}
}

func printDefaults(cmd *flag.FlagSet) {
	println(cmd.Name())
	cmd.PrintDefaults()
}

func printUsage() {
	println("Usage:")
	println()
	printDefaults(lshCmd)
	println()
	printDefaults(simCmd)
	println()
}

func parseCommand(cmd *flag.FlagSet) {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("error in parsing arguments: %v \n", err)
		printDefaults(cmd)
		os.Exit(3)
	}
}

func getShingles(source string, doKShingle bool) []string {
	textLines := getTextLines(source)

	if doKShingle {
		return lsh.KShingle(textLines, *simKShingles)
	}

	return lsh.Shingle(textLines)
}

func shingleSets(sourcesList []string, doKShingle bool) ([][]string, int) {
	fmt.Printf("\nshingling %d sources:\n", len(sourcesList))

	shingleSets := make([][]string, 0)
	var k int
	var totalSize int
	for _, s := range sourcesList {
		shingles := getShingles(s, doKShingle)
		// skip empty
		if len(shingles) == 0 {
			fmt.Printf("---> skipping %s: no shingles\n", s)
			continue
		}
		totalSize += len(shingles)
		shingleSets = append(shingleSets, shingles)
		fmt.Printf("[%d]: %s - %.150s\n", k, s, shingles[0])
		k++
	}
	return shingleSets, totalSize / len(shingleSets)
}

func doShingles(cmd *flag.FlagSet) {
	parseCommand(cmd)
	shingles := getShingles(*shingleSource, *shingleK)
	fmt.Printf("%s\n", shingles)
}

func doLSH(cmd *flag.FlagSet) {
	parseCommand(cmd)

	shingleSets, avgSize := shingleSets(toSourceList(*lshSources), false)
	if len(shingleSets) < 2 {
		fmt.Printf("nothing to compare, got %d shingle set(s)\n", len(shingleSets))
		os.Exit(0)
	}

	fmt.Printf("\naverage shingle set size is %d\n", avgSize)

	fmt.Printf("\nhashing %d sets\n", len(shingleSets))

	numHashes := *lshNumHashes
	if numHashes == 0 {
		numHashes = lsh.SuggestHashNum(avgSize)
	}

	numBands := *lshNumBands
	if numBands == 0 {
		numBands = numHashes / 5
	}

	if numBands == 0 {
		numBands = 1
	}

	fmt.Printf("\napplying %d hash functions\n", numHashes)
	signatureMatrix := lsh.Minhash(shingleSets, numHashes)

	fmt.Printf("\ndistributing into %d bands\n", numBands)
	bandBuckets := lsh.LSH(signatureMatrix, numBands)
	candidates := bandBuckets.FindCandidates()

	fmt.Printf("\nfound %d candidate pair(s)\n", len(candidates.Index))
	if len(candidates.Index) > 0 {
		fmt.Printf("%v\n", candidates.Keys())
	}
}

func doSim(cmd *flag.FlagSet) {
	parseCommand(cmd)

	texts := make([][]string, 0)
	sourceList := toSourceList(*simSources)
	for _, s := range sourceList {
		texts = append(texts, getTextLines(s))
	}

	if len(texts) != 2 {
		fmt.Println("you can compare only 2 sets")
		os.Exit(0)
	}

	fmt.Printf("similarity: %.4f\n", lsh.Jaccard(texts[0], texts[1]))
}

func parseHTML(html []string, verbose bool) []string {
	// will trim out all the tabs from text
	hizer, err := htmlizer.New([]rune{'\t'})
	if err != nil && verbose {
		fmt.Printf("error in parsing HTML lines: %v\n", err)
		return []string{}
	}

	for _, line := range html {
		err = hizer.Load(line)
		if err != nil && verbose {
			fmt.Printf("error in loading line \"%s\": %v\n", line, err)
		}
	}

	if verbose {
		fmt.Println("\nparsed HTML: ")
		fmt.Printf("%v\n\n", hizer)
	}

	return strings.Split(hizer.HumanReadable(), "\n")
}

func getTextLines(source string) []string {
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

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return parseHTML(lines, false)
	}

	return lines
}

func toSourceList(sources string) []string {
	if sources == "" {
		printUsage()
		os.Exit(4)
	}

	sourcesList := strings.Split(sources, ",")
	if len(sourcesList) < 2 {
		fmt.Println("need at least 2 documents to find candidates")
		os.Exit(0)
	}

	return sourcesList
}
