package search

import (
	"testing"
)

var doc Document = Document{
	ID:      "1",
	Content: "test match on go",
}

// func TestTokenize_SimpleSplit(t *testing.T) {
func TestAddDocument(t *testing.T) {
	invertedIndex := NewInvertedIndex()
	invertedIndex.AddDoc(doc)
	val, ok := invertedIndex.TokenIDsMap["go"]

	if !ok {
		t.Fatalf("expected true, got false")
	}

	if len(val) < 1 {
		t.Fatalf("expected 1 result, got %d", len(val))
	}

}
