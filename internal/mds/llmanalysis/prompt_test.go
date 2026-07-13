package llmanalysis

import (
	"strings"
	"testing"
	"time"

	"github.com/nahyunsama/ceactl/internal/mds/logcompressor"
)

func mustParseDay(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse("2006 Jan 2 15:04:05", value)
	if err != nil {
		t.Fatalf("failed to parse test time %q: %v", value, err)
	}

	return parsed
}

func TestBuildUserPrompt_StructuredEvents(t *testing.T) {
	result := &logcompressor.Result{
		Groups: []logcompressor.Group{
			{
				Facility: "ETHPORT",
				Mnemonic: "IF_DOWN_LINK_FAILURE",
				Iface:    "IPStorage1/6",
				Vsan:     "-",
				Severity: "5",
				Count:    4,
				First:    mustParseDay(t, "2026 Jun 1 11:34:35"),
				Last:     mustParseDay(t, "2026 Jun 1 13:55:07"),
			},
			{
				Facility: "PORT",
				Mnemonic: "IF_TRUNK_DOWN",
				Iface:    "fcip303",
				Vsan:     "22",
				Severity: "5",
				Count:    3,
				First:    mustParseDay(t, "2026 Jun 1 11:35:00"),
				Last:     mustParseDay(t, "2026 Jun 1 13:55:00"),
			},
		},
	}

	got, err := BuildUserPrompt(PromptInput{
		Device:      "GJ-IPSAN-M9K-1",
		FilterStart: mustParseDay(t, "2026 Jun 1 00:00:00"),
		FilterEnd:   mustParseDay(t, "2026 Jun 1 23:59:59"),
		Result:      result,
	})
	if err != nil {
		t.Fatalf("BuildUserPrompt returned an error: %v", err)
	}

	expected := []string{
		`device: "GJ-IPSAN-M9K-1"`,
		`filter_start: "2026-06-01 00:00:00"`,
		`filter_end: "2026-06-01 23:59:59"`,
		`event_time_min: "2026-06-01 11:34:35"`,
		`event_time_max: "2026-06-01 13:55:07"`,
		`<events count="2">`,
		`event id=1 severity="5" facility="ETHPORT"`,
		`mnemonic="IF_DOWN_LINK_FAILURE"`,
		`interface="IPStorage1/6" vsan=null observed_count=4`,
		`first="2026-06-01 11:34:35"`,
		`last="2026-06-01 13:55:07"`,
		`event id=2 severity="5" facility="PORT"`,
		`interface="fcip303" vsan="22" observed_count=3`,
	}

	for _, value := range expected {
		if !strings.Contains(got, value) {
			t.Errorf("output does not contain %q:\n%s", value, got)
		}
	}
}

func TestBuildUserPrompt_EmptyResult(t *testing.T) {
	got, err := BuildUserPrompt(PromptInput{
		Result: &logcompressor.Result{},
	})
	if err != nil {
		t.Fatalf("BuildUserPrompt returned an error: %v", err)
	}

	expected := []string{
		"device: null",
		"filter_start: null",
		"filter_end: null",
		"event_time_min: null",
		"event_time_max: null",
		`<events count="0">`,
		"</events>",
	}

	for _, value := range expected {
		if !strings.Contains(got, value) {
			t.Errorf("output does not contain %q:\n%s", value, got)
		}
	}
}

func TestBuildUserPrompt_SummarizesUnparsed(t *testing.T) {
	result := &logcompressor.Result{
		Unparsed: []string{
			"2026 Jun 1 11:35:45 switch last message repeated 2 times",
			"2026 Jun 1 11:47:50 switch last message repeated 3 times",
			"unrecognized log line",
		},
	}

	got, err := BuildUserPrompt(PromptInput{
		Result: result,
	})
	if err != nil {
		t.Fatalf("BuildUserPrompt returned an error: %v", err)
	}

	expected := []string{
		"repeat_notice_lines: 2",
		"unassigned_repeat_occurrences: 5",
		"other_unparsed_lines: 1",
	}

	for _, value := range expected {
		if !strings.Contains(got, value) {
			t.Errorf("output does not contain %q:\n%s", value, got)
		}
	}
}

func TestBuildUserPrompt_NilResult(t *testing.T) {
	_, err := BuildUserPrompt(PromptInput{})
	if err == nil {
		t.Fatal("expected an error for a nil result")
	}
}
