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

// => Got 4 matches.
```

What did we match?

```go
for _, match := range matches {
    fmt.Printf("Matched %q at position %d.\n", match.Match(), match.Pos())
}

//
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

### Building

### Searching

### Compared to Other

### Memory Usage
