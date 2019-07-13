package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	ak "github.com/anknown/ahocorasick"
	bs "github.com/bobusumisu/aho-corasick"
	cf "github.com/cloudflare/ahocorasick"
	ih "github.com/iohub/ahocorasick"
)

func timeStuff(fn func()) float64 {
	start := time.Now()
	fn()
	end := time.Now()
	return float64(end.UnixNano()-start.UnixNano()) * 1e-6
}

func totalAlloc() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.TotalAlloc) / (1024 * 1024 * 1024)
}

func iohub(patterns []string, input string) (float64, float64, int, float64) {
	m := ih.NewMatcher()

	buildTime := timeStuff(func() {
		for i, pattern := range patterns {
			m.Insert([]byte(pattern), i)
		}
		m.Compile()
	})

	seq := []byte(input)
	numMatches := 0

	searchTime := timeStuff(func() {
		resp := m.Match(seq)
		for resp.HasNext() {
			items := resp.NextMatchItem(seq)
			for range items {
				numMatches++
			}
		}
		resp.Release()
	})

	return buildTime, searchTime, numMatches, totalAlloc()
}

func anknown(patterns []string, input string) (float64, float64, int, float64) {
	runePatterns := make([][]rune, len(patterns))
	for i, pattern := range patterns {
		runePatterns[i] = []rune(pattern)
	}
	runeInput := []rune(input)

	m := new(ak.Machine)
	var matches []*ak.Term

	buildTime := timeStuff(func() {
		if err := m.Build(runePatterns); err != nil {
			panic(err)
		}
	})

	searchTime := timeStuff(func() {
		matches = m.MultiPatternSearch(runeInput, false)
	})

	return buildTime, searchTime, len(matches), totalAlloc()
}

func bobusumisu(patterns []string, input string) (float64, float64, int, float64) {
	var matches []*bs.Match
	var trie *bs.Trie

	buildTime := timeStuff(func() {
		trie = bs.NewTrieBuilder().AddStrings(patterns).Build()
	})

	searchTime := timeStuff(func() {
		matches = trie.MatchString(input)
	})

	return buildTime, searchTime, len(matches), totalAlloc()
}

func cloudflare(patterns []string, input string) (float64, float64, int, float64) {
	var m *cf.Matcher
	var matches []int
	buildTime := timeStuff(func() {
		m = cf.NewStringMatcher(patterns)
	})
	searchTime := timeStuff(func() {
		matches = m.Match([]byte(input))
	})
	return buildTime, searchTime, len(matches), totalAlloc()
}

var tests = []struct {
	name string
	fn   func([]string, string) (float64, float64, int, float64)
}{
	{"anknown   ", anknown},
	{"bobusumisu", bobusumisu},
	{"cloudflare", cloudflare},
	{"iohub     ", iohub},
}

func main() {

	f, err := os.Open("../test_data/NSF-ordlisten.cleaned.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	patterns := make([]string, 0)

	for s.Scan() {
		patterns = append(patterns, strings.TrimSpace(s.Text()))
	}

	if err := s.Err(); err != nil {
		panic(err)
	}

	inputBytes, err := ioutil.ReadFile("../test_data/Ibsen.txt")
	if err != nil {
		panic(err)
	}
	input := string(inputBytes)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "name\tpatterns\tbuild\tsearch\tmatches\talloc\t\n")
	for i := 1000; i < len(patterns); i *= 2 {
		for _, test := range tests {
			runtime.GC()
			buildTime, searchTime, numMatches, totalAlloc := test.fn(patterns[:i], input)
			fmt.Fprintf(w, "%s\t%d\t%.02fms\t%.02fms\t%d\t%.02fGiB\t\n", test.name, i, buildTime, searchTime, numMatches, totalAlloc)
		}
		fmt.Fprintf(w, "\t\t\t\t\t\t\t\n")
	}

	w.Flush()
}
