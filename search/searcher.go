package search

import (
	"cmp"
	"context"
	"slices"
)

// Searcher is the contract for anything that can answer a query with a
// ranked (or unranked) list of matching documents. The current linear
// scanner, the in-progress inverted-index + BM25 backend, and future
// vector / semantic backends all satisfy this interface, so the HTTP
// handler can stay agnostic of which implementation is plugged in.
//
// The context is threaded through so downstream implementations that
// perform I/O (LLM calls, vector DB lookups, remote index queries)
// can honour cancellation and deadlines.
type Searcher interface {
	Search(ctx context.Context, query string) []Document
}

// LinearSearcher is the naive substring-match backend. It keeps the
// documents in memory and scans them on every query. It exists as the
// baseline Searcher implementation and a reference for correctness.
type LinearSearcher struct {
	Docs []Document
}

// NewLinearSearcher builds a LinearSearcher over the given documents.
func NewLinearSearcher(docs []Document) *LinearSearcher {
	return &LinearSearcher{Docs: docs}
}

// Search satisfies the Searcher interface by delegating to the package-level
// Search function so the existing behaviour and tests are preserved.
func (s *LinearSearcher) Search(ctx context.Context, query string) []Document {
	return Search(ctx, query, s.Docs)
}

type BM25Searcher struct {
	Docs  []Document
	Index *InvertedIndex
	// Note: we keep Docs around because Search() returns []Document,
	// and the index only stores IDs. A future version could keep a
	// docID -> *Document lookup map if we want to drop the slice.
}

func NewBM25Searcher(docs []Document) *BM25Searcher {
	return &BM25Searcher{
		Docs:  docs,
		Index: BuildIndex(docs),
	}
}

type match struct {
	score    float64
	document Document
}

func (s *BM25Searcher) Search(ctx context.Context, query string) []Document {
	queryTerms := Tokenize(query)

	matches := []match{}
	for _, doc := range s.Docs {
		score := BM25Score(queryTerms, doc.ID, s.Index)
		if score > 0 {
			matches = append(matches, match{score: score, document: doc})
		}
	}

	slices.SortStableFunc(matches, func(a, b match) int {
		return cmp.Compare(b.score, a.score) // b before a = descending
	})

	results := []Document{}
	for _, m := range matches {
		results = append(results, m.document)
	}
	return results
}

var _ Searcher = (*LinearSearcher)(nil)
var _ Searcher = (*BM25Searcher)(nil)
