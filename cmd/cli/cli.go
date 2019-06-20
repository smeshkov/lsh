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

var (
	// lsh command
	lshCmd       = flag.NewFlagSet("lsh", flag.ExitOnError)
	lshSources   = lshCmd.String("s", "", "List of sources separated by comma.")
	lshNumHashes = lshCmd.Int("nh", 100, "Number of hash functions.")

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

	if doKShingle {
		return lsh.KShingle(textLines, *simKShingles)
	}

	return lsh.Shingle(textLines)
}

func shingleSets(sources string, doKShingle bool) [][]string {
	if sources == "" {
		printUsage()
		os.Exit(4)
	}

	sourcesList := strings.Split(sources, ",")
	if len(sources) < 2 {
		fmt.Println("need at least 2 documents to find candidates")
		os.Exit(0)
	}
	fmt.Printf("shingling %d sources:\n", len(sourcesList))

	shingleSets := make([][]string, 0)
	var k int
	for _, s := range sourcesList {
		shingles := getShingles(s, doKShingle)
		// skip empty
		if len(shingles) == 0 {
			fmt.Printf("---> skipping %s: no shingles\n", s)
			continue
		}
		shingleSets = append(shingleSets, shingles)
		fmt.Printf("[%d]: %s - %.150s\n", k, s, shingles[0])
		k++
	}
	return shingleSets
}

func doLSH(cmd *flag.FlagSet) {
	parseCommand(cmd)

	shingleSets := shingleSets(*lshSources, false)
	if len(shingleSets) < 2 {
		fmt.Printf("nothing to compare, got %d shingle set(s)\n", len(shingleSets))
		os.Exit(0)
	}
	fmt.Printf("\nhashing %d sets\n\n", len(shingleSets))

	signatureMatrix := lsh.Minhash(shingleSets, *lshNumHashes)
	bandBuckets := lsh.LSH(signatureMatrix, 1)
	candidates := bandBuckets.FindCandidates()

	fmt.Printf("found %d candidate pair(s)\n", len(candidates.Index))
	if len(candidates.Index) > 0 {
		fmt.Printf("%v\n", candidates.Keys())
	}
}

func doSim(cmd *flag.FlagSet) {
	parseCommand(cmd)

	shingleSets := shingleSets(*simSources, true)
	if len(shingleSets) != 2 {
		fmt.Println("you can compare only 2 sets")
		os.Exit(0)
	}

	fmt.Printf("similarity: %.4f\n", lsh.Jaccard(shingleSets[0], shingleSets[1]))
}
