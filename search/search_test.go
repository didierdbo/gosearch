package search

import "testing"

var firstDoc Document = Document{
	ID:      "1",
	Content: "test match on go",
}
var secondDoc Document = Document{
	ID:      "2",
	Content: "test another match on go",
}

func TestMatchOneDoc(t *testing.T) {
	results := Search("go", []Document{firstDoc})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].ID != firstDoc.ID {
		t.Errorf("expected ID %s, got %s", firstDoc.ID, results[0].ID)
	}
}

func TestMatchMultipleDocs(t *testing.T) {
	results := Search("go", []Document{firstDoc, secondDoc})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestNoMatch(t *testing.T) {
	results := Search("gone", []Document{firstDoc, secondDoc})

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestEmptyQuery(t *testing.T) {
	results := Search("", []Document{firstDoc, secondDoc})

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestInsensitive(t *testing.T) {
	results := Search("GO", []Document{firstDoc, secondDoc})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
