package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/didierdbo/gosearch/search"
)

func SearchHandler(docs []search.Document) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		q := r.URL.Query().Get("q")
		if q == "" {
			http.Error(w, "missing q parameter", http.StatusBadRequest)
			return
		}
		results := search.Search(q, docs)

		sr := search.SearchResult{
			Query:   q,
			Count:   len(results),
			Results: results,
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(sr)
		if err != nil {
			// http.Error(w, "encode error", http.StatusInternalServerError)
			log.Printf("encode error: %v", err)
			return
		}
	}
}
