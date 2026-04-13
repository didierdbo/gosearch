package search

import (
	"math"
	"testing"
)

func TestAddDocument(t *testing.T) {
	idx := NewInvertedIndex()
	idx.AddDoc(Document{ID: "1", Content: "test match on go"})

	tfMap, ok := idx.Postings["go"]
	if !ok {
		t.Fatalf("expected token 'go' in postings")
	}
	if tfMap["1"] != 1 {
		t.Fatalf("expected tf=1 for (go, doc1), got %d", tfMap["1"])
	}
	if idx.DocLengths["1"] != 4 {
		t.Fatalf("expected doc length 4, got %d", idx.DocLengths["1"])
	}
	if idx.TotalDocs != 1 {
		t.Fatalf("expected TotalDocs=1, got %d", idx.TotalDocs)
	}
}

func TestAddMultipleDocuments(t *testing.T) {
	idx := NewInvertedIndex()
	idx.AddDoc(Document{ID: "1", Content: "test match on go"})
	idx.AddDoc(Document{ID: "2", Content: "go is so cool, let's go learn it"})
	idx.AddDoc(Document{ID: "3", Content: "I'm going to learn a new language"})

	tfMap, exists := idx.Postings["go"]
	if !exists {
		t.Fatalf("expected token 'go' in postings")
	}
	if _, exists := idx.Postings["go"]["3"]; exists {
		t.Fatalf("doc3 should not have 'go' — no stemming expected")
	}
	if tfMap["1"] != 1 {
		t.Fatalf("expected tf=1 for (go, doc1), got %d", tfMap["1"])
	}
	if tfMap["2"] != 2 {
		t.Fatalf("expected tf=2 for (go, doc2), got %d", tfMap["2"])
	}
	if idx.DocLengths["1"] != 4 {
		t.Fatalf("expected doc1 length 4, got %d", idx.DocLengths["1"])
	}
	if idx.DocLengths["2"] != 8 {
		t.Fatalf("expected doc2 length 8, got %d", idx.DocLengths["2"])
	}
	if idx.DocLengths["3"] != 7 {
		t.Fatalf("expected doc3 length 7, got %d", idx.DocLengths["3"])
	}
	if idx.TotalDocs != 3 {
		t.Fatalf("expected TotalDocs=3, got %d", idx.TotalDocs)
	}
}

func TestBuildIndex(t *testing.T) {
	docs := []Document{
		{ID: "a", Content: "the quick brown fox"}, // 4 tokens
		{ID: "b", Content: "the lazy dog"},        // 3 tokens
		{ID: "c", Content: "quick brown quick"},   // 3 tokens
	}
	idx := BuildIndex(docs)

	if idx.TotalDocs != 3 {
		t.Fatalf("TotalDocs: want 3, got %d", idx.TotalDocs)
	}
	// avgdl = (4+3+3)/3 = 3.333...
	if math.Abs(idx.AvgDocLen-10.0/3.0) > 1e-9 {
		t.Fatalf("AvgDocLen: want ~3.333, got %f", idx.AvgDocLen)
	}
	// "quick" appears in docs a and c, with tf=1 and tf=2
	if idx.Postings["quick"]["a"] != 1 || idx.Postings["quick"]["c"] != 2 {
		t.Fatalf("quick postings wrong: %+v", idx.Postings["quick"])
	}
	// df("quick") = 2
	if len(idx.Postings["quick"]) != 2 {
		t.Fatalf("df(quick): want 2, got %d", len(idx.Postings["quick"]))
	}
}

func TestScoreTF_SaturatesWithMoreOccurrences(t *testing.T) {
	// Same docLen and avgdl, so length term is exactly 1.
	// Then scoreTF = tf*(k1+1) / (tf + k1).
	// With k1=1.2: tf=1 -> 2.2/2.2 = 1.0; tf=10 -> 22/11.2 ≈ 1.964;
	// tf=100 -> 220/101.2 ≈ 2.174. Note how slowly it grows.
	// tf=100, docLen=10 -> Note this is an impossible case, just here to validate slow increase
	k1, b := 1.2, 0.75

	s1 := scoreTF(1, 10, 10, k1, b)
	s10 := scoreTF(10, 10, 10, k1, b)
	s100 := scoreTF(100, 10, 10, k1, b)

	if math.Abs(s1-1.0) > 1e-9 {
		t.Fatalf("tf=1: want 1.0, got %f", s1)
	}
	if !(s10 > s1 && s100 > s10) {
		t.Fatalf("scores should be monotonically increasing")
	}
	if s100 > 2.2 { // must stay below the k1+1 asymptote
		t.Fatalf("tf=100: %f exceeds k1+1 asymptote", s100)
	}

	baseline := scoreTF(3, 100, 80, k1, b)
	heavilyRepeated := scoreTF(10, 10, 80, k1, b)
	rareInLongDoc := scoreTF(1, 500, 80, k1, b)
	shortDocBoost := scoreTF(1, 10, 80, k1, b)
	tfEqualsDocLen := scoreTF(5, 5, 80, k1, b)

	if !(baseline > rareInLongDoc) {
		t.Fatalf("scores should penalize long doc with rare tf")
	}
	if !(shortDocBoost > baseline) {
		t.Fatalf("short doc should score higher than baseline")
	}
	if !(heavilyRepeated < 2.2) {
		t.Fatalf("tf=10, docLen=10: %f exceeds k1+1 asymptote", heavilyRepeated)
	}
	if !(heavilyRepeated > tfEqualsDocLen) {
		t.Fatalf("high saturation should be better than lower saturation")
	}

}

func TestIDF_RareTermHigherThanCommonTerm(t *testing.T) {
	rare := idf(1, 1000)          // in 1 doc out of 1000
	common := idf(500, 1000)      // in half the corpus
	everywhere := idf(1000, 1000) // in every doc

	if !(rare > common && common > everywhere) {
		t.Fatalf("IDF ordering wrong: rare=%f common=%f everywhere=%f",
			rare, common, everywhere)
	}
	if everywhere < 0 {
		t.Fatalf("smoothed IDF must not go negative, got %f", everywhere)
	}
}

func TestBM25Score_HandComputed(t *testing.T) {
	docs := []Document{
		{ID: "a", Content: "the quick brown fox"},
		{ID: "b", Content: "the lazy dog"},
		{ID: "c", Content: "quick brown quick"},
	}
	idx := BuildIndex(docs)

	// Query: "quick"
	//   N = 3, df(quick) = 2, so idf = ln((3-2+0.5)/(2+0.5) + 1) = ln(1.6)
	//   For doc c: tf=2, |D|=3, avgdl=10/3
	//     length_norm = 1 - 0.75 + 0.75 * 3 / (10/3) = 0.25 + 0.675 = 0.925
	//     denom = 2 + 1.2 * 0.925 = 3.11
	//     numer = 2 * 2.2 = 4.4
	//     tf_component = 4.4 / 3.11 ≈ 1.4148
	//     score = ln(1.6) * 1.4148 ≈ 0.470 * 1.4148 ≈ 0.6650
	//
	// Compute this by hand, paste the expected value below, and
	// let the test enforce your understanding.
	got := BM25Score([]string{"quick"}, "c", idx)
	want := 0.6650
	if math.Abs(got-want) > 1e-3 {
		t.Fatalf("BM25Score: want ~%f, got %f", want, got)
	}
}
