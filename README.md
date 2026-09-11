# name-forge

A command-line tool that prints random, human-readable names like
`brave-falcon-42`. Useful for anything that needs a throwaway identifier a
person can actually read and say out loud: container names, feature
branches, test fixtures, temp directories.

## Usage

Build and run:

```
go build -o namegen .
./namegen
brave-falcon-742
```

Generate several at once:

```
./namegen -n 5
amber-canyon-19
tidy-glacier-501
keen-marsh-8
...
```

Change the separator or drop the numeric suffix:

```
./namegen -sep _ -number=false
gentle_orchard
```

### Custom wordlists

By default the tool uses a small built-in list of adjectives and nouns. You
can point it at your own files instead, one word per line:

```
./namegen -adjectives ./my-adjectives.txt -nouns ./my-nouns.txt -n 3
```

This is where the tool earns its keep: `-adjectives`/`-nouns` files are
streamed line by line and never loaded into memory in full. A 50 MB
wordlist and a 5 KB one cost the same amount of memory to sample from -
only the words that actually end up in the output are held onto, via
reservoir sampling. You can safely point this at a dictionary-sized file.

Either flag also accepts `-` to read the list from stdin instead of a file,
which is handy for piping in a filtered or generated list:

```
grep -v '^#' my-nouns.txt | ./namegen -nouns - -n 3
```

`-adjectives` and `-nouns` can't both be `-` in the same run, since stdin
can only be consumed once.

## Flags

| Flag          | Default | Meaning                                        |
|---------------|---------|-------------------------------------------------|
| `-n`          | `1`     | number of names to generate                     |
| `-sep`        | `-`     | separator between words                         |
| `-adjectives` | (built-in list) | path to a custom adjective file, one per line |
| `-nouns`      | (built-in list) | path to a custom noun file, one per line      |
| `-number`     | `true`  | append a random number suffix                   |
| `-max`        | `1000`  | exclusive upper bound for the number suffix     |
| `-seed`       | `0`     | random seed; `0` derives one from the clock      |

## Status

Early skeleton. Word selection and the streaming sampler work and are
covered by unit tests; see the repo for what's planned next.
