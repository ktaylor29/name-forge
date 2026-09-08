package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeWordlist(t *testing.T, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "words.txt")
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		t.Fatalf("writing test wordlist: %v", err)
	}
	return path
}

func TestSampleLinesReturnsKLines(t *testing.T) {
	path := writeWordlist(t, []string{"a", "b", "c", "d", "e"})
	rng := rand.New(rand.NewSource(1))

	got, err := sampleLines(path, 3, rng)
	if err != nil {
		t.Fatalf("sampleLines() error = %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("sampleLines() returned %d lines, want 3", len(got))
	}
}

func TestSampleLinesOnlyReturnsWordsFromFile(t *testing.T) {
	source := []string{"a", "b", "c", "d", "e"}
	path := writeWordlist(t, source)
	rng := rand.New(rand.NewSource(1))

	got, err := sampleLines(path, 3, rng)
	if err != nil {
		t.Fatalf("sampleLines() error = %v", err)
	}
	valid := map[string]bool{}
	for _, w := range source {
		valid[w] = true
	}
	for _, w := range got {
		if !valid[w] {
			t.Errorf("sampleLines() returned %q, not present in source file", w)
		}
	}
}

func TestSampleLinesSkipsBlankLines(t *testing.T) {
	path := writeWordlist(t, []string{"a", "", "b", "", "", "c"})
	rng := rand.New(rand.NewSource(1))

	got, err := sampleLines(path, 3, rng)
	if err != nil {
		t.Fatalf("sampleLines() error = %v", err)
	}
	for _, w := range got {
		if w == "" {
			t.Errorf("sampleLines() returned a blank line")
		}
	}
}

func TestSampleLinesToppedUpWhenFileHasFewerWordsThanK(t *testing.T) {
	path := writeWordlist(t, []string{"a", "b"})
	rng := rand.New(rand.NewSource(1))

	got, err := sampleLines(path, 5, rng)
	if err != nil {
		t.Fatalf("sampleLines() error = %v", err)
	}
	if len(got) != 5 {
		t.Fatalf("sampleLines() returned %d lines, want 5", len(got))
	}
	for _, w := range got {
		if w != "a" && w != "b" {
			t.Errorf("sampleLines() returned %q, want a or b", w)
		}
	}
}

func TestSampleLinesErrorsOnMissingFile(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	if _, err := sampleLines(filepath.Join(t.TempDir(), "missing.txt"), 1, rng); err == nil {
		t.Error("sampleLines() error = nil, want error for missing file")
	}
}

func TestSampleLinesErrorsOnNoWords(t *testing.T) {
	path := writeWordlist(t, []string{"", "", ""})
	rng := rand.New(rand.NewSource(1))
	if _, err := sampleLines(path, 1, rng); err == nil {
		t.Error("sampleLines() error = nil, want error for file with no usable words")
	}
}

func TestSampleLinesErrorsOnInvalidK(t *testing.T) {
	path := writeWordlist(t, []string{"a"})
	rng := rand.New(rand.NewSource(1))
	if _, err := sampleLines(path, 0, rng); err == nil {
		t.Error("sampleLines() error = nil, want error for k < 1")
	}
}

func TestSampleLinesUniformCoverageAcrossFile(t *testing.T) {
	source := []string{"a", "b", "c", "d", "e"}
	path := writeWordlist(t, source)
	rng := rand.New(rand.NewSource(7))

	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		got, err := sampleLines(path, 1, rng)
		if err != nil {
			t.Fatalf("sampleLines() error = %v", err)
		}
		seen[got[0]] = true
	}
	for _, w := range source {
		if !seen[w] {
			t.Errorf("word %q was never sampled across 200 draws", w)
		}
	}
}
