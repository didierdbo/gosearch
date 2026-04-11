package search

import "testing"

func assertStringSlices(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length: got %d, want %d", len(got), len(want))
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("%d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestTokenize_SimpleSplit(t *testing.T) {
	input := "hello world"
	expected := []string{"hello", "world"}
	result := Tokenize(input)
	assertStringSlices(t, result, expected)
}

func TestTokenize_Lowercase(t *testing.T) {
	input := "Hello WORLD"
	expected := []string{"hello", "world"}
	result := Tokenize(input)
	assertStringSlices(t, result, expected)
}

func TestTokenize_PunctuationRemoved(t *testing.T) {
	input := "hello, world!"
	expected := []string{"hello", "world"}
	result := Tokenize(input)
	assertStringSlices(t, result, expected)
}
func TestTokenize_MultipleSpaces(t *testing.T) {
	input := "hello  world"
	expected := []string{"hello", "world"}
	result := Tokenize(input)
	assertStringSlices(t, result, expected)
}
func TestTokenize_EmptyString(t *testing.T) {
	input := ""
	expected := []string{}
	result := Tokenize(input)
	assertStringSlices(t, result, expected)
}
