package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/smeshkov/lsh"
)

func main() {
	reader := readStdIN()
	lines, err := readLines(reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", lsh.Shingle(lines))
}

func readStdIN() *bufio.Reader {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(fmt.Sprintf("error in reading from STDIN: %v", err))
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		panic("unsupported mode")
	}
	return bufio.NewReader(os.Stdin)
}

// readLines provides slice of strings from input split by white space.
func readLines(reader io.Reader) ([]string, error) {
	defer close(reader)

	scanner := bufio.NewScanner(reader)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// Close ...
func close(reader io.Reader) {
	if closer, ok := reader.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
	}
}
