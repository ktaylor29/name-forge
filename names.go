package main

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"
)

// Supported values for the -format flag.
const (
	FormatKebab = "kebab"
	FormatCamel = "camel"
	FormatTitle = "title"
)

// defaultAdjectives and defaultNouns back the tool when the caller doesn't
// supply their own wordlists. Kept short on purpose - anyone who wants more
// variety or a different flavor (Latin plant names, Norse gods, whatever)
// can point -adjectives/-nouns at their own file instead of us maintaining
// an ever-growing embedded dictionary.
var defaultAdjectives = []string{
	"amber", "brave", "calm", "dusty", "eager", "faded", "gentle", "hollow",
	"idle", "jagged", "keen", "lively", "muted", "narrow", "olive", "plain",
	"quiet", "rusty", "solid", "tidy", "urban", "vivid", "warm", "young",
	"blunt", "crisp", "dense", "even", "fond", "grim",
}

var defaultNouns = []string{
	"badger", "canyon", "delta", "ember", "falcon", "glacier", "harbor",
	"island", "jungle", "kettle", "lagoon", "meadow", "nebula", "orchard",
	"pebble", "quarry", "ridge", "summit", "thicket", "valley", "willow",
	"cinder", "drift", "quartz", "fjord", "grove", "hollow", "inlet",
	"marsh", "reef",
}

// buildName picks one adjective and one noun at random, joins them
// according to format, and optionally appends a random number so repeat
// runs don't collide as often (e.g. for use as container or branch names).
//
// camelCase has no separator by definition, so sep is ignored between
// words (and before the number suffix) when format is FormatCamel.
func buildName(rng *rand.Rand, adjectives, nouns []string, sep string, withNumber bool, maxNumber int, format string) string {
	adj := adjectives[rng.Intn(len(adjectives))]
	noun := nouns[rng.Intn(len(nouns))]

	var name string
	switch format {
	case FormatCamel:
		name = adj + capitalize(noun)
	case FormatTitle:
		name = strings.Join([]string{capitalize(adj), capitalize(noun)}, sep)
	default:
		name = strings.Join([]string{adj, noun}, sep)
	}

	if withNumber {
		if format == FormatCamel {
			name = fmt.Sprintf("%s%d", name, rng.Intn(maxNumber))
		} else {
			name = fmt.Sprintf("%s%s%d", name, sep, rng.Intn(maxNumber))
		}
	}

	return name
}

// capitalize upper-cases the first rune of s and leaves the rest alone, so
// words with non-ASCII first letters still work.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
