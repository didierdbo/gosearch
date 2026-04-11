package search

import "context"

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

var _ Searcher = (*LinearSearcher)(nil)
