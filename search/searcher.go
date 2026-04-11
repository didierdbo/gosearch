package search

// Searcher is the contract for anything that can answer a query with a
// ranked (or unranked) list of matching documents. The current linear
// scanner, the in-progress inverted-index + BM25 backend, and future
// vector / semantic backends all satisfy this interface, so the HTTP
// handler can stay agnostic of which implementation is plugged in.
type Searcher interface {
	Search(query string) []Document
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
func (s *LinearSearcher) Search(query string) []Document {
	return Search(query, s.Docs)
}
