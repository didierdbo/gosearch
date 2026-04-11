package search

type InvertedIndex struct {
	TokenIDsMap map[string][]string `json:"token_ids"`
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		TokenIDsMap: map[string][]string{},
	}
}

func (index *InvertedIndex) AddDoc(doc Document) error {
	result := Tokenize(doc.Content)
	for _, token := range result {
		ids, ok := index.TokenIDsMap[token]
		if !ok {
			ids = []string{}
		}
		ids = append(ids, doc.ID)
		index.TokenIDsMap[token] = ids
	}

	return nil
}
