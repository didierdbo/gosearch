package search

type InvertedIndex struct {
	Postings   map[string]map[string]int `json:"postings"`
	DocLengths map[string]int            `json:"doc_lengths"`
	TotalDocs  int                       `json:"total_docs"`
	AvgDocLen  float64                   `json:"avg_doc_len"`
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Postings:   make(map[string]map[string]int),
		DocLengths: make(map[string]int),
		TotalDocs:  0,
		AvgDocLen:  0.0,
	}
}

func BuildIndex(docs []Document) *InvertedIndex {
	idx := NewInvertedIndex()
	for _, doc := range docs {
		idx.AddDoc(doc)
	}
	if idx.TotalDocs > 0 {
		idx.AvgDocLen = float64(sum(idx.DocLengths)) / float64(idx.TotalDocs)
	}

	return idx
}

func sum(lengths map[string]int) int {
	total := 0
	for _, length := range lengths {
		total += length
	}
	return total
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
