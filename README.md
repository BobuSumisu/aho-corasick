# Aho-Corasick

> forked from [BobuSumisu/aho-corasick](https://github.com/BobuSumisu/aho-corasick)
>
> Fix the bug of missing pattern field in the return value (Match) of MatchFirst function


---


[![Build Status](https://travis-ci.com/BobuSumisu/aho-corasick.svg?token=eGRFn5xdQ7p9yby3GVvc&branch=master)](https://travis-ci.com/BobuSumisu/aho-corasick)
![Go Version](https://img.shields.io/github/go-mod/go-version/BobuSumisu/aho-corasick)
![Latest Tag](https://img.shields.io/github/v/tag/BobuSumisu/aho-corasick)

Implementation of the Aho-Corasick string-search algorithm in Go.

Licensed under MIT License.

## Details

This implementation does not use a [Double-Array Trie](https://linux.thai.net/~thep/datrie/datrie.html) as in my
[implementation](https://github.com/BobuSumisu/go-ahocorasick) from a couple of years back.

This reduces the build time drastically, but at the cost of higher memory consumption.

The search time is still fast, and comparable to other Go implementations I have found on github that claims to be fast
(see [performance](#Performance)).

## Documentation

Can be found at [godoc.org](https://godoc.org/github.com/BobuSumisu/aho-corasick).

## Example Usage

Use a `TrieBuilder` to build a `Trie`:

```go
trie := NewTrieBuilder().
    AddStrings([]string{"or", "amet"}).
    Build()
```

Then go and match something interesting:

```go
matches := trie.MatchString("Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
fmt.Printf("Got %d matches.\n", len(matches))

// => Got 3 matches.
```

What did we match?

```go
for _, match := range matches {
    fmt.Printf("Matched pattern %d %q at position %d.\n", match.Match(),
        match.Pattern(), match.Pos())
}

// => Matched pattern 0 "or" at position 1.
// => Matched pattern 0 "or" at position 15.
// => Matched patterh 1 "amet" at position 22.
```

## Building

You can easily load patterns from file:

```go
builder := NewTrieBuilder()
builder.LoadPatterns("patterns.txt")
builder.LoadStrings("strings.txt")
```

Both functions expects a text file with one pattern per line. `LoadPatterns` expects the pattern to
be in hexadecimal form.

## Storing

Use `Encode` to store a `Trie` in gzip compressed binary format:

```go
f, err := os.Create("trie.gz")
err := Encode(f, trie)
```

And `Decode` to load it from binary format:

```go
f, err := os.Open("trie.gz")
trie, err := Decode(f)
```

## Performance

Some simple benchmarking on my machine (Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz, 32 GiB RAM).

Build and search time grows quite linearly with regards to number of patterns and input text length.

### Building

    BenchmarkTrieBuild/100-4                    7886            154786 ns/op
    BenchmarkTrieBuild/1000-4                    739           1647419 ns/op
    BenchmarkTrieBuild/10000-4                    91          13331713 ns/op
    BenchmarkTrieBuild/100000-4                    9         123886615 ns/op


### Searching

    BenchmarkMatchIbsen/100-4                1471089               819 ns/op
    BenchmarkMatchIbsen/1000-4                202288              5667 ns/op
    BenchmarkMatchIbsen/10000-4                19957             59680 ns/op
    BenchmarkMatchIbsen/100000-4                2012            595086 ns/op

### Compared to Other Implementation

See
[aho-corasick-benchmark](https://github.com/Bobusumisu/aho-corasick-benchmark).

### Memory Usage

As mentioned, the memory consumption will be quite high compared to a double-array trie
implementation. Especially during the build phase (which currently contains a lot of object
allocations).
