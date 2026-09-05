package main

import (
	"fmt"
	"math/rand"
	"strings"
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

// buildName picks one adjective and one noun at random, joins them with sep,
// and optionally appends a random number so repeat runs don't collide as
// often (e.g. for use as container or branch names).
func buildName(rng *rand.Rand, adjectives, nouns []string, sep string, withNumber bool, maxNumber int) string {
	name := strings.Join([]string{
		adjectives[rng.Intn(len(adjectives))],
		nouns[rng.Intn(len(nouns))],
	}, sep)

	if withNumber {
		name = fmt.Sprintf("%s%s%d", name, sep, rng.Intn(maxNumber))
	}

	return name
}
