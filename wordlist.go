package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
)

// stdinPath is the conventional "read from stdin instead of a file" marker,
// same as cat, tar, and friends use.
const stdinPath = "-"

// sampleLines picks k lines out of the file at path using reservoir
// sampling, so memory use stays at O(k) no matter how large the file is.
// A naive approach (read all lines, shuffle, slice) would need to hold the
// entire wordlist in memory, which defeats the point of letting people point
// this at a big dictionary file. path may be "-" to read from stdin.
func sampleLines(path string, k int, rng *rand.Rand) ([]string, error) {
	r := os.Stdin
	if path != stdinPath {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}

	reservoir, err := sampleLinesFromReader(r, k, rng)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	return reservoir, nil
}

// sampleLinesFromReader does the actual reservoir sampling over r. Split out
// from sampleLines so both file paths and stdin can share it.
func sampleLinesFromReader(r io.Reader, k int, rng *rand.Rand) ([]string, error) {
	if k < 1 {
		return nil, fmt.Errorf("sampleLinesFromReader: k must be at least 1")
	}

	reservoir := make([]string, 0, k)
	scanner := bufio.NewScanner(r)
	seen := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		seen++
		if len(reservoir) < k {
			reservoir = append(reservoir, line)
			continue
		}
		// classic Algorithm R: replace an existing slot with probability k/seen
		j := rng.Intn(seen)
		if j < k {
			reservoir[j] = line
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if seen == 0 {
		return nil, fmt.Errorf("no words found (blank lines are skipped)")
	}

	// file had fewer usable lines than requested names; top up by sampling
	// with replacement from what we actually found
	for len(reservoir) < k {
		reservoir = append(reservoir, reservoir[rng.Intn(len(reservoir))])
	}

	return reservoir, nil
}
