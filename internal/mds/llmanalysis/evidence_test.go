package llmanalysis

import (
	"reflect"
	"testing"
)

func TestReferencedEventIDs_ReturnsUniqueValidIDsInFirstSeenOrder(t *testing.T) {
	reply := "Evidence: E3, E1, E3, E8, and E99"

	got := ReferencedEventIDs(reply, 8)
	want := []int{3, 1, 8}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ReferencedEventIDs() = %v, want %v", got, want)
	}
}

func TestReferencedEventIDs_ExcludesRangeReferences(t *testing.T) {
	reply := "Ranges E1-E9, E10 through E20, and E21 to E23; individual E4 and E7."

	got := ReferencedEventIDs(reply, 23)
	want := []int{4, 7}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ReferencedEventIDs() = %v, want %v", got, want)
	}
}

func TestReferencedEventIDs_EmptyWhenNoValidReferenceExists(t *testing.T) {
	if got := ReferencedEventIDs("No cited events", 5); len(got) != 0 {
		t.Fatalf("ReferencedEventIDs() = %v, want no IDs", got)
	}

	if got := ReferencedEventIDs("Evidence: E1", 0); len(got) != 0 {
		t.Fatalf("ReferencedEventIDs() = %v, want no IDs", got)
	}
}
