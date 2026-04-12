package search

type InvertedIndex struct {
	Postings     map[string]map[string]int `json:"postings"`
	DocLengths   map[string]int            `json:"doc_lengths"`
	TotalDocs    int                       `json:"total_docs"`
	AvgDocLength float64                   `json:"avg_doc_len"`
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Postings:     make(map[string]map[string]int),
		DocLengths:   make(map[string]int),
		TotalDocs:    0,
		AvgDocLength: 0.0,
	}
}

func (index *InvertedIndex) AddDoc(doc Document) error {
	result := Tokenize(doc.Content)
	for _, token := range result {
		if _, ok := index.Postings[token]; !ok {
			index.Postings[token] = make(map[string]int)
		}
		index.Postings[token][doc.ID]++
	}
	if _, ok := index.DocLengths[doc.ID]; !ok {
		index.TotalDocs++
	}
	index.DocLengths[doc.ID] = len(result)

	return nil
}
