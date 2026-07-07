package logcompressor

import (
	"strings"
	"testing"
	"time"
)

func mustParseDay(t *testing.T, s string) time.Time {
	t.Helper()
	tm, err := time.Parse("2006 Jan 2 15:04:05", s)
	if err != nil {
		t.Fatalf("failed to parse test time %q: %v", s, err)
	}
	return tm
}

func TestAnalyze_GroupsIdenticalEvents(t *testing.T) {
	log := strings.Join([]string{
		"2024 Jan 15 10:23:45 switch1 %PORT-5-IF_DOWN: Interface fc1/1, VSAN 100 is down",
		"2024 Jan 15 10:24:10 switch1 %PORT-5-IF_DOWN: Interface fc1/1, VSAN 100 is down",
	}, "\n")

	result, err := Analyze(strings.NewReader(log), time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(result.Groups))
	}

	g := result.Groups[0]
	if g.Count != 2 {
		t.Errorf("Count = %d, want 2", g.Count)
	}
	if g.Facility != "PORT" || g.Mnemonic != "IF_DOWN" || g.Iface != "fc1/1" || g.Vsan != "100" {
		t.Errorf("unexpected group fields: %+v", g)
	}

	wantFirst := mustParseDay(t, "2024 Jan 15 10:23:45")
	wantLast := mustParseDay(t, "2024 Jan 15 10:24:10")
	if !g.First.Equal(wantFirst) {
		t.Errorf("First = %v, want %v", g.First, wantFirst)
	}
	if !g.Last.Equal(wantLast) {
		t.Errorf("Last = %v, want %v", g.Last, wantLast)
	}
}

func TestAnalyze_DifferentInterfacesCreateSeparateGroups(t *testing.T) {
	log := strings.Join([]string{
		"2024 Jan 15 10:23:45 switch1 %PORT-5-IF_DOWN: Interface fc1/1, VSAN 100 is down",
		"2024 Jan 15 10:23:46 switch1 %PORT-5-IF_DOWN: Interface fc1/2, VSAN 100 is down",
	}, "\n")

	result, err := Analyze(strings.NewReader(log), time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Groups) != 2 {
		t.Fatalf("got %d groups, want 2", len(result.Groups))
	}
}

func TestAnalyze_UnparsedLines(t *testing.T) {
	log := strings.Join([]string{
		"this line has no timestamp at all",
		"2024 Jan 15 10:23:45 switch1 no mnemonic here",
		"",
		"2024 Jan 15 10:23:46 switch1 %PORT-5-IF_DOWN: Interface fc1/1, VSAN 100 is down",
	}, "\n")

	result, err := Analyze(strings.NewReader(log), time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(result.Groups))
	}
	if len(result.Unparsed) != 2 {
		t.Fatalf("got %d unparsed lines, want 2 (blank line must be skipped entirely): %v", len(result.Unparsed), result.Unparsed)
	}
}

func TestAnalyze_FromToFiltersOutOfRangeLines(t *testing.T) {
	log := strings.Join([]string{
		"2024 Jan 15 09:00:00 switch1 %PORT-5-IF_DOWN: Interface fc1/1, VSAN 100 is down",
		"2024 Jan 15 10:00:00 switch1 %PORT-5-IF_UP: Interface fc1/1, VSAN 100 is up",
		"2024 Jan 15 11:00:00 switch1 %PORT-5-IF_DOWN: Interface fc1/2, VSAN 200 is down",
	}, "\n")

	from := mustParseDay(t, "2024 Jan 15 09:30:00")
	to := mustParseDay(t, "2024 Jan 15 10:30:00")

	result, err := Analyze(strings.NewReader(log), from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Groups) != 1 {
		t.Fatalf("got %d groups, want 1 (only the 10:00:00 line falls in range)", len(result.Groups))
	}
	if result.Groups[0].Mnemonic != "IF_UP" {
		t.Errorf("got mnemonic %q, want IF_UP", result.Groups[0].Mnemonic)
	}
	if len(result.Unparsed) != 0 {
		t.Errorf("got %d unparsed, want 0 (out-of-range lines are dropped, not marked unparsed)", len(result.Unparsed))
	}
}

func TestAnalyze_GroupsOrderedByFirstOccurrence(t *testing.T) {
	// Group with iface fc1/2 is scanned first, but its timestamp is later than
	// the fc1/1 group's. The final result should still sort by First ascending.
	log := strings.Join([]string{
		"2024 Jan 15 10:00:00 switch1 %PORT-5-IF_DOWN: Interface fc1/2, VSAN 200 is down",
		"2024 Jan 15 09:00:00 switch1 %PORT-5-IF_DOWN: Interface fc1/1, VSAN 100 is down",
	}, "\n")

	result, err := Analyze(strings.NewReader(log), time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Groups) != 2 {
		t.Fatalf("got %d groups, want 2", len(result.Groups))
	}
	if result.Groups[0].Iface != "fc1/1" {
		t.Errorf("Groups[0].Iface = %q, want fc1/1 (earlier timestamp should sort first)", result.Groups[0].Iface)
	}
	if result.Groups[1].Iface != "fc1/2" {
		t.Errorf("Groups[1].Iface = %q, want fc1/2", result.Groups[1].Iface)
	}
}
