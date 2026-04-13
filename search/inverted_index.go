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

// scoreTF returns the BM25 term-frequency component for a single
// (term, document) pair. It is the fraction on the right-hand side
// of the BM25 sum, without the IDF factor.
//
// Parameters:
//
//	tf     : term frequency of the term in the document
//	docLen : length of the document in tokens
//	avgdl  : average document length across the corpus
//	k1, b  : BM25 tuning parameters
func scoreTF(tf, docLen int, avgdl, k1, b float64) float64 {
	// TODO: write this using the formula from Section 1.
	// Hint: the numerator is tf*(k1+1), the denominator involves
	// tf + k1 * (1 - b + b * docLen/avgdl). Cast ints to float64.
	num := float64(tf) * (k1 + 1)
	den := float64(tf) + k1*(1-b+b*(float64(docLen)/avgdl))

	return num / den
}
