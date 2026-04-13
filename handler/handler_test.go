package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didierdbo/gosearch/search"
)

var firstDoc search.Document = search.Document{
	ID:      "1",
	Content: "test match on go",
}
var secondDoc search.Document = search.Document{
	ID:      "2",
	Content: "test another match on go",
}

func TestValidRequest(t *testing.T) {
	request := httptest.NewRequest("GET", "/search?q=go", nil)
	recorder := httptest.NewRecorder()
	searcher := search.NewLinearSearcher([]search.Document{firstDoc, secondDoc})
	SearchHandler(searcher)(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestMissingQuery(t *testing.T) {
	request := httptest.NewRequest("GET", "/search", nil)
	recorder := httptest.NewRecorder()
	searcher := search.NewLinearSearcher([]search.Document{firstDoc, secondDoc})
	SearchHandler(searcher)(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestBadHttpMethod(t *testing.T) {
	request := httptest.NewRequest("POST", "/search?q=go", nil)
	recorder := httptest.NewRecorder()
	searcher := search.NewLinearSearcher([]search.Document{firstDoc, secondDoc})
	SearchHandler(searcher)(recorder, request)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", recorder.Code)
	}
}

func TestSearchHandler_BM25EndToEnd(t *testing.T) {
	docs := []search.Document{
		{ID: "a", Content: "the quick brown fox"},
		{ID: "c", Content: "quick brown quick"},
	}
	s := search.NewBM25Searcher(docs)
	h := SearchHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/search?q=quick+brown", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var results search.SearchResult

	out, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if err := json.Unmarshal(out, &results); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if results.Count != 2 {
		t.Fatalf("want 2 results, got %d", results.Count)
	}

	if results.Results[0].ID != "c" {
		t.Fatalf("want doc c first, got %s", results.Results[0].ID)
	}

	if results.Results[1].ID != "a" {
		t.Fatalf("want doc a second, got %s", results.Results[1].ID)
	}
}
