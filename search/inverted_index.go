package search

import (
	"math"
)

const (
	BM25K1 = 1.2
	BM25B  = 0.75
)

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
	num := float64(tf) * (k1 + 1)
	den := float64(tf) + k1*(1-b+b*(float64(docLen)/avgdl))
	return num / den
}

// idf computes the smoothed inverse document frequency for a term.
//
// Parameters:
//   df : number of documents containing the term (document frequency)
//   N  : total number of documents in the corpus
//
// Uses the BM25+ / Lucene smoothing:
//   idf = ln( (N - df + 0.5) / (df + 0.5) + 1 )
func idf(df, N int) float64 {
	dfFloat := float64(df)
	nFloat := float64(N)
	return math.Log(((nFloat - dfFloat + 0.5) / (dfFloat + 0.5)) + 1)
}

// BM25Score computes the BM25 relevance score of a document for a
// tokenized query. Returns 0 if the document does not contain any
// of the query terms.
//
// queryTerms : already tokenized (use the same Tokenize() as indexing!)
// docID      : the document to score
// idx        : the fully built inverted index
func BM25Score(queryTerms []string, docID string, idx *InvertedIndex) float64 {
	var score float64
	for _, term := range queryTerms {
		df := len(idx.Postings[term])
		if df == 0 {
			continue
		}
		tf := idx.Postings[term][docID]
		if tf == 0 {
			continue
		}
		termIDF := idf(df, idx.TotalDocs)
		tfComponent := scoreTF(tf, idx.DocLengths[docID], idx.AvgDocLen, BM25K1, BM25B)
		score += termIDF * tfComponent
	}
	return score
}
