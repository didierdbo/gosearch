package search

import (
	"context"
	"testing"
)

var firstDoc Document = Document{
	ID:      "1",
	Content: "test match on go",
}
var secondDoc Document = Document{
	ID:      "2",
	Content: "test another match on go",
}

func TestMatchOneDoc(t *testing.T) {
	results := Search(context.Background(), "go", []Document{firstDoc})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].ID != firstDoc.ID {
		t.Errorf("expected ID %s, got %s", firstDoc.ID, results[0].ID)
	}
}

func TestMatchMultipleDocs(t *testing.T) {
	results := Search(context.Background(), "go", []Document{firstDoc, secondDoc})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestNoMatch(t *testing.T) {
	results := Search(context.Background(), "gone", []Document{firstDoc, secondDoc})

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestEmptyQuery(t *testing.T) {
	results := Search(context.Background(), "", []Document{firstDoc, secondDoc})

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestInsensitive(t *testing.T) {
	results := Search(context.Background(), "GO", []Document{firstDoc, secondDoc})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestBM25Searcher_RanksByRelevance(t *testing.T) {
	docs := []Document{
		{ID: "a", Content: "the quick brown fox"},
		{ID: "b", Content: "the lazy dog"},
		{ID: "c", Content: "quick brown quick"},
	}
	s := NewBM25Searcher(docs)
	results := s.Search(context.Background(), "quick brown")

	if len(results) != 2 {
		t.Fatalf("want 2 results, got %d", len(results))
	}
	// doc c mentions "quick" twice and "brown" once, doc a mentions each once.
	// c should rank first.
	if results[0].ID != "c" {
		t.Fatalf("want doc c first, got %s", results[0].ID)
	}
}
