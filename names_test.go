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

	got := buildName(rng, adjectives, nouns, "-", false, 1000, FormatKebab)
	if got != "brave-falcon" {
		t.Errorf("buildName() = %q, want %q", got, "brave-falcon")
	}
}

func TestBuildNameUsesSeparator(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	adjectives := []string{"brave"}
	nouns := []string{"falcon"}

	got := buildName(rng, adjectives, nouns, "_", false, 1000, FormatKebab)
	if got != "brave_falcon" {
		t.Errorf("buildName() = %q, want %q", got, "brave_falcon")
	}
}

func TestBuildNameWithoutNumberOmitsSuffix(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	got := buildName(rng, []string{"brave"}, []string{"falcon"}, "-", false, 1000, FormatKebab)
	if strings.Count(got, "-") != 1 {
		t.Errorf("buildName() = %q, want exactly one separator with -number=false", got)
	}
}

func TestBuildNameWithNumberAppendsBoundedSuffix(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 200; i++ {
		got := buildName(rng, []string{"brave"}, []string{"falcon"}, "-", true, 5, FormatKebab)
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

func TestBuildNameCamelFormatHasNoSeparators(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	got := buildName(rng, []string{"brave"}, []string{"falcon"}, "-", true, 5, FormatCamel)
	if !strings.HasPrefix(got, "braveFalcon") {
		t.Errorf("buildName() = %q, want prefix %q", got, "braveFalcon")
	}
	if strings.Contains(got, "-") {
		t.Errorf("buildName() = %q, want no separators in camel format", got)
	}
}

func TestBuildNameTitleFormatCapitalizesEachWord(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	got := buildName(rng, []string{"brave"}, []string{"falcon"}, "-", false, 1000, FormatTitle)
	if got != "Brave-Falcon" {
		t.Errorf("buildName() = %q, want %q", got, "Brave-Falcon")
	}
}

func TestCombosMultipliesListsAndNumberRange(t *testing.T) {
	adjectives := []string{"amber", "brave", "calm"}
	nouns := []string{"badger", "canyon"}

	if got, want := combos(adjectives, nouns, false, 1000), 6; got != want {
		t.Errorf("combos() = %d, want %d", got, want)
	}
	if got, want := combos(adjectives, nouns, true, 10), 60; got != want {
		t.Errorf("combos() = %d, want %d", got, want)
	}
}

func TestBuildUniqueNameNeverRepeats(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	adjectives := []string{"amber", "brave", "calm"}
	nouns := []string{"badger", "canyon"}
	total := combos(adjectives, nouns, false, 1000)

	seen := map[string]bool{}
	for i := 0; i < total; i++ {
		name, err := buildUniqueName(rng, adjectives, nouns, "-", false, 1000, FormatKebab, seen)
		if err != nil {
			t.Fatalf("buildUniqueName() error = %v on draw %d", err, i)
		}
		if _, ok := seen[name]; !ok {
			t.Fatalf("buildUniqueName() returned %q but didn't record it in seen", name)
		}
	}
	if len(seen) != total {
		t.Errorf("buildUniqueName() produced %d distinct names, want %d", len(seen), total)
	}
}

func TestBuildUniqueNameErrorsWhenSpaceExhausted(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	adjectives := []string{"brave"}
	nouns := []string{"falcon"}
	seen := map[string]bool{"brave-falcon": true}

	if _, err := buildUniqueName(rng, adjectives, nouns, "-", false, 1000, FormatKebab, seen); err == nil {
		t.Error("buildUniqueName() error = nil, want error when only possible name is already seen")
	}
}

func TestBuildNamePicksFromWholeList(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	adjectives := []string{"amber", "brave", "calm"}
	nouns := []string{"badger", "canyon"}

	seenAdj := map[string]bool{}
	seenNoun := map[string]bool{}
	for i := 0; i < 200; i++ {
		got := buildName(rng, adjectives, nouns, "-", false, 1000, FormatKebab)
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
