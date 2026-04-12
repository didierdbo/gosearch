package search

import (
	"testing"
)

func TestAddDocument(t *testing.T) {
	idx := NewInvertedIndex()
	idx.AddDoc(Document{ID: "1", Content: "test match on go"})

	tfMap, ok := idx.Postings["go"]
	if !ok {
		t.Fatalf("expected token 'go' in postings")
	}
	if tfMap["1"] != 1 {
		t.Fatalf("expected tf=1 for (go, doc1), got %d", tfMap["1"])
	}
	if idx.DocLengths["1"] != 4 {
		t.Fatalf("expected doc length 4, got %d", idx.DocLengths["1"])
	}
	if idx.TotalDocs != 1 {
		t.Fatalf("expected TotalDocs=1, got %d", idx.TotalDocs)
	}
}
