package ahocorasick

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func testTrie(trie *Trie) error {
	matches := trie.MatchString("Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
	expected := []*Match{
		newMatchString(1, "or"),
		newMatchString(15, "or"),
		newMatchString(22, "amet"),
	}

	if len(expected) != len(matches) {
		return fmt.Errorf("expected %d matches, got %d\n", len(expected), len(matches))
	}

	for i := range matches {
		if !MatchEqual(expected[i], matches[i]) {
			return fmt.Errorf("expected %v, got %v\n", expected[i], matches[i])
		}
	}

	return nil
}

func TestEncodingAndDecoding(t *testing.T) {
	trie := NewTrieBuilder().AddStrings([]string{"or", "amet"}).Build()

	if err := testTrie(trie); err != nil {
		t.Error(err)
	}

	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.Encode(trie); err != nil {
		t.Error(err)
	}

	dec := NewDecoder(&buf)

	decodedTrie, err := dec.Decode()
	if err != nil {
		t.Error(err)
	}

	if err := testTrie(decodedTrie); err != nil {
		t.Error(err)
	}
}

func TestReadAndWriteTrie(t *testing.T) {
	patterns, err := readPatterns("test_data/NSF-ordlisten.cleaned.uniq.txt")
	if err != nil {
		t.Fatal(err)
	}

	trie := NewTrieBuilder().AddStrings(patterns[:10000]).Build()

	defer os.Remove("test.trie")

	if err := WriteTrie(trie, "test.trie"); err != nil {
		t.Fatal(err)
	}

	decodedTrie, err := ReadTrie("test.trie")
	if err != nil {
		t.Fatal(err)
	}

	matches := decodedTrie.MatchString("abasien")

	fmt.Println(matches)

	if len(matches) != 3 {
		t.Errorf("expected 3 matches, got %d", len(matches))
	}
}
