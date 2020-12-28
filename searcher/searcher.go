package searcher

import (
	"bytes"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"sort"
	"strings"
)

// Searcher searcher struct to handle with text file
type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

// Load load file to Searcher struct
func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(bytes.ToUpper(dat))
	return nil
}

// Search returns the list of texts with searched word.
//
// Main features are:
// - Case insensitive search
// - Sort result
// - Start of each result in the beginning of the line making results more readable
//
// What not done:
// - Not remove duplicated result in the same line because prioritized only one word in each result
// - Result not limited, if look for "lord" it can delay too much. The result was not limited because
// the it if necessary can be cached and return by pages
// - Markdown answer not added to not added not builtin package, but it can be added later
func (s *Searcher) Search(query string) []string {
	lineDelimiter := []byte("\n")
	const maxChar = 200
	const delimiterString = "***"

	lq := len(query)
	query = strings.ToUpper(query)
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	// find all lines
	idxs2 := s.SuffixArray.Lookup(lineDelimiter, -1)
	// sort lookup results to better visualization
	sort.Ints(idxs2)
	sort.Ints(idxs)
	results := []string{}
	for _, idx := range idxs {
		// look for start-end line
		r := binarySearch(idx, idxs2)
		id2 := r[0]
		idxRel := idx - id2
		snip := s.CompleteWorks[id2 : id2+maxChar]
		// added strings before and after word found to improve visualization
		// the string added is *** that can be converted to bold in markdown
		snip = snip[:idxRel] + delimiterString + snip[idxRel:idxRel+lq] + delimiterString + snip[idxRel+lq:]
		snips := strings.Split(snip, " ")
		// trim spaces, ignore last word if it is not a full word and add "..."
		snip = strings.TrimSpace(strings.Join(snips[:len(snips)-1], " ")) + "..."
		results = append(results, string(snip))
	}
	return results
}

// binarySearch recursive binary search algorithm to find delimiters and
// returns a slice with the start and end of line
func binarySearch(t int, p []int) []int {
	l := (len(p) >> 1)
	if len(p) == 3 {
		if t > p[1] {
			return p[1:3]
		}
		return p[:2]
	}
	if l <= 1 {
		return p
	}
	if t < p[l] {
		if t > p[l-1] {
			return p[l-1 : l+1]
		}
		return binarySearch(t, p[:l])
	}
	return binarySearch(t, p[l:])
}
