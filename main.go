package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Document struct {
	ID      string
	Content string
}

var docs = []Document{
	{ID: "1", Content: "sur le search"},
	{ID: "2", Content: "Go"},
	{ID: "3", Content: "Python"},
}

func main() {
	http.HandleFunc("/search", searchHandler)
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func search(query string, docs []Document) []Document {
	results := []Document{}
	for _, doc := range docs {
		if strings.Contains(strings.ToLower(doc.Content), strings.ToLower(query)) {
			results = append(results, doc)
		}
	}
	return results
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "missing q parameter", http.StatusBadRequest)
		return
	}
	results := search(q, docs)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		// http.Error(w, "encode error", http.StatusInternalServerError)
		log.Printf("encode error: %v", err)
		return
	}
}
