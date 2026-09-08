package main

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func TestBuildNameJoinsAdjectiveAndNoun(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	adjectives := []string{"brave"}
	nouns := []string{"falcon"}

	got := buildName(rng, adjectives, nouns, "-", false, 1000)
	if got != "brave-falcon" {
		t.Errorf("buildName() = %q, want %q", got, "brave-falcon")
	}
}

func TestBuildNameUsesSeparator(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	adjectives := []string{"brave"}
	nouns := []string{"falcon"}

	got := buildName(rng, adjectives, nouns, "_", false, 1000)
	if got != "brave_falcon" {
		t.Errorf("buildName() = %q, want %q", got, "brave_falcon")
	}
}

func TestBuildNameWithoutNumberOmitsSuffix(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	got := buildName(rng, []string{"brave"}, []string{"falcon"}, "-", false, 1000)
	if strings.Count(got, "-") != 1 {
		t.Errorf("buildName() = %q, want exactly one separator with -number=false", got)
	}
}

func TestBuildNameWithNumberAppendsBoundedSuffix(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 200; i++ {
		got := buildName(rng, []string{"brave"}, []string{"falcon"}, "-", true, 5)
		parts := strings.Split(got, "-")
		if len(parts) != 3 {
			t.Fatalf("buildName() = %q, want 3 parts separated by -", got)
		}
		if parts[0] != "brave" || parts[1] != "falcon" {
			t.Fatalf("buildName() = %q, want brave-falcon-N", got)
		}
		n, err := strconv.Atoi(parts[2])
		if err != nil {
			t.Fatalf("buildName() suffix %q is not a number: %v", parts[2], err)
		}
		if n < 0 || n >= 5 {
			t.Errorf("buildName() suffix %d out of range [0,5)", n)
		}
	}
}

func TestBuildNamePicksFromWholeList(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	adjectives := []string{"amber", "brave", "calm"}
	nouns := []string{"badger", "canyon"}

	seenAdj := map[string]bool{}
	seenNoun := map[string]bool{}
	for i := 0; i < 200; i++ {
		got := buildName(rng, adjectives, nouns, "-", false, 1000)
		parts := strings.SplitN(got, "-", 2)
		seenAdj[parts[0]] = true
		seenNoun[parts[1]] = true
	}
	for _, a := range adjectives {
		if !seenAdj[a] {
			t.Errorf("adjective %q was never picked across 200 draws", a)
		}
	}
	for _, n := range nouns {
		if !seenNoun[n] {
			t.Errorf("noun %q was never picked across 200 draws", n)
		}
	}
}
