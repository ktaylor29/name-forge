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

Change the casing with `-format`:

```
./namegen -format camel
braveFalcon742

./namegen -format title
Brave-Falcon-742
```

`camel` ignores `-sep` entirely, since camelCase has no separators by
definition. `title` and the default `kebab` both still honor `-sep`.

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

Avoid repeats within a single run with `-unique`:

```
./namegen -n 5 -unique
```

`-unique` rejects the run up front if `-n` asks for more names than the
wordlists and `-max` can actually produce (e.g. two one-word lists with
`-number=false` can only ever make one name).

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
| `-format`     | `kebab` | output casing: `kebab`, `camel`, or `title`      |
| `-unique`     | `false` | never print the same name twice in one run      |
| `-version`    | `false` | print the version and exit                      |

## Releases

Releases are built with [goreleaser](https://goreleaser.com) from
`.goreleaser.yaml`. It runs the tests, then builds static binaries for
linux, darwin and windows on amd64 and arm64, and writes archives plus a
checksums file to `dist/`. To check the config and build locally without
publishing:

```
goreleaser release --snapshot --clean
```

To cut a real release, tag the commit and run `goreleaser release --clean`
with `GITHUB_TOKEN` set. Binaries report the tag from `namegen -version`;
a plain `go build` reports `dev`.

## Status

Early skeleton. Word selection and the streaming sampler work and are
covered by unit tests; see the repo for what's planned next.
