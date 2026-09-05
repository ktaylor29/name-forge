// Command namegen prints memorable random names, e.g. "brave-falcon-42",
// suitable for container names, feature branches, or test fixtures.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	count := flag.Int("n", 1, "number of names to generate")
	sep := flag.String("sep", "-", "separator between words")
	adjPath := flag.String("adjectives", "", "path to a custom adjective list, one word per line")
	nounPath := flag.String("nouns", "", "path to a custom noun list, one word per line")
	withNumber := flag.Bool("number", true, "append a random number suffix")
	maxNumber := flag.Int("max", 1000, "exclusive upper bound for the number suffix")
	seed := flag.Int64("seed", 0, "random seed; 0 derives a seed from the current time")
	flag.Parse()

	if *count < 1 {
		fmt.Fprintln(os.Stderr, "namegen: -n must be at least 1")
		os.Exit(1)
	}
	if *maxNumber < 1 {
		fmt.Fprintln(os.Stderr, "namegen: -max must be at least 1")
		os.Exit(1)
	}

	s := *seed
	if s == 0 {
		s = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(s))

	adjectives := defaultAdjectives
	nouns := defaultNouns

	if *adjPath != "" {
		list, err := sampleLines(*adjPath, *count, rng)
		if err != nil {
			fmt.Fprintf(os.Stderr, "namegen: %v\n", err)
			os.Exit(1)
		}
		adjectives = list
	}
	if *nounPath != "" {
		list, err := sampleLines(*nounPath, *count, rng)
		if err != nil {
			fmt.Fprintf(os.Stderr, "namegen: %v\n", err)
			os.Exit(1)
		}
		nouns = list
	}

	for i := 0; i < *count; i++ {
		fmt.Println(buildName(rng, adjectives, nouns, *sep, *withNumber, *maxNumber))
	}
}
