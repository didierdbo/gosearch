package handler

import (
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
