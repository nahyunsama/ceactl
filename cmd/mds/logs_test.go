package mds

import (
	"testing"
	"time"
)

func TestParseDayStart(t *testing.T) {
	t.Run("empty string returns zero time with no error", func(t *testing.T) {
		got, err := parseDayStart("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !got.IsZero() {
			t.Errorf("got %v, want zero time", got)
		}
	})

	t.Run("valid date returns midnight of that day", func(t *testing.T) {
		got, err := parseDayStart("20240115")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)
		if !got.Equal(want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("invalid format returns an error", func(t *testing.T) {
		_, err := parseDayStart("2024-01-15")
		if err == nil {
			t.Fatal("expected an error for malformed date, got nil")
		}
	})
}

func TestParseDayEnd(t *testing.T) {
	t.Run("empty string returns zero time with no error", func(t *testing.T) {
		got, err := parseDayEnd("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !got.IsZero() {
			t.Errorf("got %v, want zero time", got)
		}
	})

	t.Run("valid date returns the last second of that day", func(t *testing.T) {
		got, err := parseDayEnd("20240115")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := time.Date(2024, time.January, 15, 23, 59, 59, 0, time.UTC)
		if !got.Equal(want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("invalid format returns an error", func(t *testing.T) {
		_, err := parseDayEnd("2024-01-15")
		if err == nil {
			t.Fatal("expected an error for malformed date, got nil")
		}
	})
}
