package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/didierdbo/gosearch/handler"
	"github.com/didierdbo/gosearch/search"
)

func main() {
	docs, err := loadDocuments("documents.json")
	if err != nil {
		log.Fatal(err)
	}

	searcher := search.NewLinearSearcher(docs)
	http.HandleFunc("/search", handler.SearchHandler(searcher))
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadDocuments(path string) ([]search.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	docs := []search.Document{}
	if err := json.Unmarshal(data, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}
