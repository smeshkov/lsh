# LSH (Locality-Sensitive Hashing)

[![Build Status](https://travis-ci.com/smeshkov/lsh.svg?branch=master)](https://travis-ci.com/smeshkov/lsh)
[![Go Report Card](https://goreportcard.com/badge/github.com/smeshkov/lsh)](https://goreportcard.com/report/github.com/smeshkov/lsh)
[![Coverage](https://codecov.io/gh/smeshkov/lsh/branch/master/graph/badge.svg)](https://codecov.io/gh/smeshkov/lsh)
[![GoDoc](https://godoc.org/github.com/smeshkov/lsh?status.svg)](https://godoc.org/github.com/smeshkov/lsh)

Finding similar items with "Locality-Sensitive Hashing". Inspired by Chapter 3: "Finding Similar Items" of "Mining Massive Datasets" book by Jure Leskovec, Anand Rajaraman and Jeff Ullman.

# Usage

See [examples](https://github.com/smeshkov/lsh/tree/master/examples) for use cases.

Steps in code:

1. `#Shingle` - tokenize
2. `#Minhash` - signature matrix
3. `#LSH` - candidate pairs
3. `#Jaccard` - for jaccard similarity of candidate pairs

in CLI: `./lsh lsh -s <comma_separated_URLs>`. For example:

```bash
./lsh lsh -s https://stackoverflow.com,https://stackoverflow.com
shingling 2 sources:
[0]: https://stackoverflow.com - more stack exchange
[1]: https://stackoverflow.com - more stack exchange

hashing 2 sets

found 1 candidate pair(s)
[0_1]
```

This means that sources 0 and 1 ended up as a candidate pair, 
therefore they are suggested for similarity test.


Then `./lsh sim -s <two_comma_separated_URLs>`. For example:

```bash
./lsh sim -s https://stackoverflow.com,https://stackoverflow.com
shingling 2 sources:
[0]: https://stackoverflow.com - more stack exchange
[1]: https://stackoverflow.com - more stack exchange
similarity: 1.0000
```

# Check-list

### Similarity

- [x] Jaccard

### Performance tests

TODO.

## Changelog

See [CHANGELOG.md](https://raw.githubusercontent.com/smeshkov/lsh/master/CHANGELOG.md)

## License

Released under the [Apache License 2.0](https://raw.githubusercontent.com/smeshkov/lsh/master/LICENSE).
