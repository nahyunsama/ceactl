package llmanalysis

import (
	"strings"
	"testing"
	"time"

	"github.com/nahyunsama/ceactl/internal/mds/logcompressor"
)

func mustParseDay(t *testing.T, s string) time.Time {
	t.Helper()
	tm, err := time.Parse("2006 Jan 2 15:04:05", s)
	if err != nil {
		t.Fatalf("failed to parse test time %q: %v", s, err)
	}
	return tm
}

func TestBuildUserPrompt_IncludesDeviceAndRange(t *testing.T) {
	result := &logcompressor.Result{
		Groups: []logcompressor.Group{
			{
				Facility: "PORT", Mnemonic: "IF_DOWN", Iface: "fc1/1", Vsan: "100", Severity: "5", Count: 3,
				First: mustParseDay(t, "2026 Jun 1 11:34:35"),
				Last:  mustParseDay(t, "2026 Jun 1 13:55:07"),
			},
		},
	}

	got, err := BuildUserPrompt("mds-lab-1", result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, "device (mds-lab-1)") {
		t.Errorf("output missing device name:\n%s", got)
	}
	if !strings.Contains(got, "2026 Jun 1, 11:34:35–13:55:07") {
		t.Errorf("output missing same-day time range:\n%s", got)
	}
	if !strings.Contains(got, "[sev5] PORT-IF_DOWN") {
		t.Errorf("output missing group table row:\n%s", got)
	}
}

func TestBuildUserPrompt_RangeSpansMultipleDays(t *testing.T) {
	result := &logcompressor.Result{
		Groups: []logcompressor.Group{
			{
				Facility: "PORT", Mnemonic: "IF_DOWN", Iface: "fc1/1", Vsan: "100", Severity: "5", Count: 1,
				First: mustParseDay(t, "2026 Jun 1 23:55:00"),
				Last:  mustParseDay(t, "2026 Jun 2 00:05:00"),
			},
		},
	}

	got, err := BuildUserPrompt("mds-lab-1", result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, "2026 Jun 1, 23:55:00–2026 Jun 2, 00:05:00") {
		t.Errorf("output missing multi-day time range:\n%s", got)
	}
}

func TestBuildUserPrompt_NoGroupsHasUnknownRange(t *testing.T) {
	got, err := BuildUserPrompt("mds-lab-1", &logcompressor.Result{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, "an unknown time range") {
		t.Errorf("output missing unknown-range fallback:\n%s", got)
	}
}

func TestBuildUserPrompt_UnparsedNoteSumsRepeatCounts(t *testing.T) {
	result := &logcompressor.Result{
		Unparsed: []string{
			"2026 Jun 1 11:35:45 GJ-IPSAN-M9K-1 last message repeated 2 times",
			"2026 Jun 1 11:47:50 GJ-IPSAN-M9K-1 last message repeated 3 times",
		},
	}

	got, err := BuildUserPrompt("mds-lab-1", result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, `Note: 2 lines in the original log were device-side "last message repeated N times" compressions (totaling ~5 additional occurrences)`) {
		t.Errorf("output missing unparsed note:\n%s", got)
	}
}

func TestBuildUserPrompt_NoUnparsedOmitsNote(t *testing.T) {
	got, err := BuildUserPrompt("mds-lab-1", &logcompressor.Result{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(got, "Note:") {
		t.Errorf("output should omit the unparsed note when there are no unparsed lines:\n%s", got)
	}
}
