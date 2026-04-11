package search

import (
	"strings"
)

type Document struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type SearchResult struct {
	Query   string     `json:"query"`
	Count   int        `json:"count"`
	Results []Document `json:"results"`
}

func Search(query string, docs []Document) []Document {
	results := []Document{}
	if query == "" {
		return results
	}
	for _, doc := range docs {
		if strings.Contains(strings.ToLower(doc.Content), strings.ToLower(query)) {
			results = append(results, doc)
		}
	}
	return results
}
