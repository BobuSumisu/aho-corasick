# Aho-Corasick

> forked from [BobuSumisu/aho-corasick](https://github.com/BobuSumisu/aho-corasick)
>
> Fix the bug of missing pattern field in the return value (Match) of MatchFirst function


---

Implementation of the Aho-Corasick string-search algorithm in Go.

Licensed under MIT License.

## Documentation

Can be found at [pkg.go.dev](https://pkg.go.dev/github.com/zeewell/aho-corasick).

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
