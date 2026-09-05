package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

// sampleLines picks k lines out of the file at path using reservoir
// sampling, so memory use stays at O(k) no matter how large the file is.
// A naive approach (read all lines, shuffle, slice) would need to hold the
// entire wordlist in memory, which defeats the point of letting people point
// this at a big dictionary file.
func sampleLines(path string, k int, rng *rand.Rand) ([]string, error) {
	if k < 1 {
		return nil, fmt.Errorf("sampleLines: k must be at least 1")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reservoir := make([]string, 0, k)
	scanner := bufio.NewScanner(f)
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
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	if seen == 0 {
		return nil, fmt.Errorf("%s: no words found (blank lines are skipped)", path)
	}

	// file had fewer usable lines than requested names; top up by sampling
	// with replacement from what we actually found
	for len(reservoir) < k {
		reservoir = append(reservoir, reservoir[rng.Intn(len(reservoir))])
	}

	return reservoir, nil
}
