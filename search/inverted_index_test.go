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

func TestAddMultipleDocuments(t *testing.T) {
	idx := NewInvertedIndex()
	idx.AddDoc(Document{ID: "1", Content: "test match on go"})
	idx.AddDoc(Document{ID: "2", Content: "go is so cool, let's go learn it"})
	idx.AddDoc(Document{ID: "3", Content: "I'm going to learn a new language"})

	tfMap, exists := idx.Postings["go"]
	if !exists {
		t.Fatalf("expected token 'go' in postings")
	}
	if _, exists := idx.Postings["go"]["3"]; exists {
		t.Fatalf("doc3 should not have 'go' — no stemming expected")
	}
	if tfMap["1"] != 1 {
		t.Fatalf("expected tf=1 for (go, doc1), got %d", tfMap["1"])
	}
	if tfMap["2"] != 2 {
		t.Fatalf("expected tf=2 for (go, doc2), got %d", tfMap["2"])
	}
	if idx.DocLengths["1"] != 4 {
		t.Fatalf("expected doc1 length 4, got %d", idx.DocLengths["1"])
	}
	if idx.DocLengths["2"] != 8 {
		t.Fatalf("expected doc2 length 8, got %d", idx.DocLengths["2"])
	}
	if idx.DocLengths["3"] != 7 {
		t.Fatalf("expected doc3 length 7, got %d", idx.DocLengths["3"])
	}
	if idx.TotalDocs != 3 {
		t.Fatalf("expected TotalDocs=3, got %d", idx.TotalDocs)
	}
}
