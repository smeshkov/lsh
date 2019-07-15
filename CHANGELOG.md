# Changelog
All notable changes to this project will be documented in this file.

## 0.7.0
 - added `shingle` command to `lsh` CLI;
 - added poor man's suggestion for number of hash functions based on average sie of shingle sets;
 - bumped `github.com/zoomio/inout` to `0.6.0`.

## 0.6.0
 - changes to CLI: use raw text lines of provided sources instead of tags from given sources.

## 0.5.0
 - use K-shingling instead of stop-word based shingling for similarity comparison in CLI.

## 0.4.0
 - added subcommands to CLI: `lsh` for candidate pairs and `sim` for similarity of candidate pair.

## 0.3.0
 - added `#Jaccard` for finding Jaccard similarity between two sets.

## 0.2.2
 - first release, provides candidate pairs via pipeline: `#Shingle` -> `#Minhash` -> `#LSH`.