# Aho-Corasick

[![Build Status](https://travis-ci.com/BobuSumisu/aho-corasick.svg?token=eGRFn5xdQ7p9yby3GVvc&branch=master)](https://travis-ci.com/BobuSumisu/aho-corasick)

Implementation of the Aho-Corasick string-search algorithm in Go.

Licensed under MIT License.

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
    fmt.Printf("Matched %q at position %d.\n", match.Match(), match.Pos())
}

// => Matched "or" at position 1.
// => Matched "or" at position 15.
// => Matched "amet" at position 22.
```

## Building

You can easily load patterns from file:

```go
builder := NewTrieBuilder()
builder.LoadPatterns("patterns.txt")
builder.LoadStrings("strings.txt")
```

`LoadPatterns` expects patterns as hexadecimal strings, and both functions expects one
pattern per line.

## Saving and Loading

A `Trie` can be saved to file and loaded, to avoid having to build a trie multiple times:

```go
if err := trie.Save("my.trie"); err != nil {
    log.Fatal(err)
}
```

Loading it:

```go
trie, err := LoadTrie("my.trie")
if err != nil {
    log.Fatal(err)
}
```

## Performance

Some simple benchmarking on my machine (Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz, 32 GiB RAM).

Build and search time grows quite linearly with regards to number of patterns and input text length.

### Building

    BenchmarkTrieBuild/100-12                    10000              0.1460 ms/op
    BenchmarkTrieBuild/1000-12                    1000              2.1643 ms/op
    BenchmarkTrieBuild/10000-12                    100             14.3305 ms/op
    BenchmarkTrieBuild/100000-12                    10            131.2442 ms/op

### Searching

    BenchmarkMatchIbsen/100-12                 2000000              0.0006 ms/op
    BenchmarkMatchIbsen/1000-12                 300000              0.0042 ms/op
    BenchmarkMatchIbsen/10000-12                 30000              0.0436 ms/op
    BenchmarkMatchIbsen/100000-12                 3000              0.4310 ms/op

### Compared to Others

Inspired by [anknown](https://github.com/anknown/ahocorasick) I also wanted to check how my implementation compared
to other Aho-Corasick implementations in Go.

I created a simple [benchmark](./benchmark/main.go) and ran it on my laptop.

![benchmark plot](./benchmark/benchmark.png)
