package ahocorasick

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestReadme(t *testing.T) {
	trie := NewTrieBuilder().AddStrings([]string{"or", "amet"}).Build()
	matches := trie.MatchString("Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
	expected := []*Match{
		newMatchString(1, 0, "or"),
		newMatchString(15, 0, "or"),
		newMatchString(22, 1, "amet"),
	}

	if len(expected) != len(matches) {
		t.Errorf("expected %d matches, got %d\n", len(expected), len(matches))
	}

	for i := range matches {
		if !MatchEqual(expected[i], matches[i]) {
			t.Errorf("expected %v, got %v\n", expected[i], matches[i])
		}
	}
}

func TestTrie(t *testing.T) {
	cases := []struct {
		name     string
		patterns []string
		input    string
		expected []*Match
	}{
		{
			"Wikipedia",
			[]string{"a", "ab", "bab", "bc", "bca", "c", "caa"},
			"abccab",
			[]*Match{
				newMatchString(0, 0, "a"),
				newMatchString(0, 1, "ab"),
				newMatchString(1, 3, "bc"),
				newMatchString(2, 5, "c"),
				newMatchString(3, 5, "c"),
				newMatchString(4, 0, "a"),
				newMatchString(4, 1, "ab"),
			},
		},
		{
			"Prefix",
			[]string{"Aho-Corasick", "Aho-Cora", "Aho", "A"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(0, 3, "A"),
				newMatchString(0, 2, "Aho"),
				newMatchString(0, 1, "Aho-Cora"),
				newMatchString(0, 0, "Aho-Corasick"),
			},
		},
		{
			"Suffix",
			[]string{"Aho-Corasick", "Corasick", "sick", "k"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(0, 0, "Aho-Corasick"),
				newMatchString(4, 1, "Corasick"),
				newMatchString(8, 2, "sick"),
				newMatchString(11, 3, "k"),
			},
		},
		{
			"Infix",
			[]string{"Aho-Corasick", "ho-Corasi", "o-Co", "-"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(3, 3, "-"),
				newMatchString(2, 2, "o-Co"),
				newMatchString(1, 1, "ho-Corasi"),
				newMatchString(0, 0, "Aho-Corasick"),
			},
		},
		{
			"Overlap",
			[]string{"Aho-Co", "ho-Cora", "o-Coras", "-Corasick"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(0, 0, "Aho-Co"),
				newMatchString(1, 1, "ho-Cora"),
				newMatchString(2, 2, "o-Coras"),
				newMatchString(3, 3, "-Corasick"),
			},
		},
		{
			"Adjacent",
			[]string{"Ah", "o-Co", "ras", "ick"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(0, 0, "Ah"),
				newMatchString(2, 1, "o-Co"),
				newMatchString(6, 2, "ras"),
				newMatchString(9, 3, "ick"),
			},
		},
		{
			"SingleSymbol",
			[]string{"o"},
			"Aho-Corasick",
			[]*Match{
				newMatchString(2, 0, "o"),
				newMatchString(5, 0, "o"),
			},
		},
		{
			"NoMatch",
			[]string{"Gazorpazopfield", "Knuth", "O"},
			"Aho-Corasick",
			[]*Match{},
		},
		{
			"Zeroes",
			[]string{"\x00\x00"},
			"\x00\x00Aho\x00\x00-\x00\x00Corasick\x00\x00",
			[]*Match{
				newMatchString(0, 0, "\x00\x00"),
				newMatchString(5, 0, "\x00\x00"),
				newMatchString(8, 0, "\x00\x00"),
				newMatchString(18, 0, "\x00\x00"),
			},
		},
		{
			"Alphabetsize",
			[]string{"\xff\xff"},
			"\xff\xffAho\xfe\xfe-\xff\xffCorasick\xff\xff\xff",
			[]*Match{
				newMatchString(0, 0, "\xff\xff"),
				newMatchString(8, 0, "\xff\xff"),
				newMatchString(18, 0, "\xff\xff"),
				newMatchString(19, 0, "\xff\xff"),
			},
		},
	}

	for _, c := range cases {
		tr := NewTrieBuilder().AddStrings(c.patterns).Build()
		matches := tr.MatchString(c.input)

		if len(matches) != len(c.expected) {
			t.Errorf("%s: expected %d matches, got %d", c.name, len(c.expected), len(matches))
			continue
		}

		for i := range matches {
			if !MatchEqual(matches[i], c.expected[i]) {
				t.Errorf("%s: expected %v, got %v", c.name, c.expected[i], matches[i])
			}
		}
	}
}

func TestMatchFirst(t *testing.T) {
	ibsen, err := ioutil.ReadFile("./test_data/Ibsen.txt")
	if err != nil {
		t.Error(err)
	}
	tr := NewTrieBuilder().AddString("Hedvig").Build()
	match := tr.MatchFirst(ibsen)
	expected := newMatchString(937, 0, "Hedvig")
	if !MatchEqual(expected, match) {
		t.Errorf("expected %v, got %v\n", expected, match)
	}
}

func TestHedvig(t *testing.T) {
	ibsen, err := ioutil.ReadFile("./test_data/Ibsen.txt")
	if err != nil {
		t.Error(err)
	}
	matches := NewTrieBuilder().AddString("Hedvig").Build().Match(ibsen)
	if len(matches) != 134 {
		fmt.Printf("expected to find 134 Hedvig's, got %d\n", len(matches))
	}
}

func BenchmarkTrieBuild(b *testing.B) {
	patterns, err := readPatterns("./test_data/NSF-ordlisten.cleaned.txt")
	if err != nil {
		b.Error(err)
	}

	b.Run("100", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			NewTrieBuilder().AddStrings(patterns[:100]).Build()
		}
	})
	b.Run("1000", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			NewTrieBuilder().AddStrings(patterns[:1000]).Build()
		}
	})
	b.Run("10000", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			NewTrieBuilder().AddStrings(patterns[:10000]).Build()
		}
	})
	b.Run("100000", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			NewTrieBuilder().AddStrings(patterns[:100000]).Build()
		}
	})
}

func BenchmarkMatchIbsen(b *testing.B) {
	patterns, err := readPatterns("./test_data/NSF-ordlisten.cleaned.txt")
	if err != nil {
		b.Error(err)
	}

	ibsen, err := ioutil.ReadFile("./test_data/Ibsen.txt")
	if err != nil {
		b.Error(err)
	}

	trie := NewTrieBuilder().AddStrings(patterns[:10000]).Build()

	b.Run("100", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			trie.Match(ibsen[:100])
		}
	})
	b.Run("1000", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			trie.Match(ibsen[:1000])
		}
	})
	b.Run("10000", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			trie.Match(ibsen[:10000])
		}
	})
	b.Run("100000", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			trie.Match(ibsen[:100000])
		}
	})
}

func readPatterns(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	patterns := make([]string, 0)

	for s.Scan() {
		patterns = append(patterns, strings.TrimSpace(s.Text()))
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}
