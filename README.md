# LSH (Work-In-Progress)

[![Build Status](https://travis-ci.com/smeshkov/lsh.svg?branch=master)](https://travis-ci.com/smeshkov/lsh)
[![Go Report Card](https://goreportcard.com/badge/github.com/smeshkov/lsh)](https://goreportcard.com/report/github.com/smeshkov/lsh)

Finding similar items with "Locality-Sensitive Hashing". Inspired by Chapter 3: "Finding Similar Items" of "Mining Massive Datasets" book by Jure Leskovec, Anand Rajaraman and Jeff Ullman.

# Usage

Steps in code:

1. #Shingle
2. #Minhash
3. #LSH

OR with CLI: `lsh -s <comma_separated_URLs>`. For example:

```bash
lsh -s https://stackoverflow.com,https://stackoverflow.com
shingling 2 sources:
[0]: https://stackoverflow.com - more stack exchange
[1]: https://stackoverflow.com - more stack exchange

hashing 2 sets

found 1 candidate pair(s)
[0_1]
```

# Check-list

### Hash functions

- [x] Figure if there is a good programmatic approach to providing hash functions, e.g. define 5-10 patterns, then apply randomization
- [x] Increase amount to 100 functions

### LSH

- [x] Write tests
- [ ] Write better tests

### Similarity

TODO.

### Performance tests

TODO.

## Changelog

See [CHANGELOG.md](https://raw.githubusercontent.com/smeshkov/lsh/master/CHANGELOG.md)

## License

Released under the [Apache License 2.0](https://raw.githubusercontent.com/smeshkov/lsh/master/LICENSE).