package logcompressor

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWriteReport_HeaderCounts(t *testing.T) {
	result := &Result{
		Groups: []Group{
			{
				Facility: "PORT",
				Mnemonic: "IF_DOWN",
				Iface:    "fc1/1",
				Vsan:     "100",
				Severity: "5",
				Count:    3,
				First:    time.Date(2024, time.January, 15, 10, 23, 45, 0, time.UTC),
				Last:     time.Date(2024, time.January, 15, 10, 25, 0, 0, time.UTC),
			},
		},
		Unparsed: []string{"unparsed line 1", "unparsed line 2"},
	}

	var buf bytes.Buffer
	if err := result.WriteReport(&buf, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "그룹 1개 (미분류 2줄)") {
		t.Errorf("output missing header counts:\n%s", out)
	}
	if !strings.Contains(out, "[E1 sev5] PORT-IF_DOWN") {
		t.Errorf("output missing group facility/mnemonic:\n%s", out)
	}
	if !strings.Contains(out, "iface=fc1/1") || !strings.Contains(out, "vsan=100") {
		t.Errorf("output missing iface/vsan fields:\n%s", out)
	}
	if !strings.Contains(out, "3회") {
		t.Errorf("output missing count:\n%s", out)
	}
	if !strings.Contains(out, "10:23:45 ~ 10:25:00") {
		t.Errorf("output missing time span for differing First/Last:\n%s", out)
	}
	if !strings.Contains(out, "미분류 줄 (2개)") {
		t.Errorf("output missing unparsed section header:\n%s", out)
	}
	if !strings.Contains(out, "unparsed line 1") || !strings.Contains(out, "unparsed line 2") {
		t.Errorf("output missing unparsed lines:\n%s", out)
	}
}

func TestWriteReport_SingleTimestampHasNoRangeSeparator(t *testing.T) {
	sameTime := time.Date(2024, time.January, 15, 10, 23, 45, 0, time.UTC)
	result := &Result{
		Groups: []Group{
			{Facility: "PORT", Mnemonic: "IF_DOWN", Iface: "fc1/1", Vsan: "100", Severity: "5", Count: 1, First: sameTime, Last: sameTime},
		},
	}

	var buf bytes.Buffer
	if err := result.WriteReport(&buf, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if strings.Contains(out, "~") {
		t.Errorf("output should not show a range when First equals Last:\n%s", out)
	}
	if !strings.Contains(out, "(10:23:45)") {
		t.Errorf("output missing single timestamp:\n%s", out)
	}
}

func TestWriteReport_LimitsUnparsedLines(t *testing.T) {
	result := &Result{
		Unparsed: []string{"line 1", "line 2", "line 3"},
	}

	var buf bytes.Buffer
	if err := result.WriteReport(&buf, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "line 1") {
		t.Errorf("output should include the first unparsed line within the limit:\n%s", out)
	}
	if strings.Contains(out, "line 2") || strings.Contains(out, "line 3") {
		t.Errorf("output should not include lines beyond the limit:\n%s", out)
	}
	if !strings.Contains(out, "외 2줄") {
		t.Errorf("output missing truncation notice for remaining lines:\n%s", out)
	}
}

func TestWriteReport_TrimsUnparsedLineWhitespace(t *testing.T) {
	result := &Result{
		Unparsed: []string{"  padded line  "},
	}

	var buf bytes.Buffer
	if err := result.WriteReport(&buf, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "padded line") {
		t.Errorf("output missing trimmed line:\n%s", out)
	}
	if strings.Contains(out, "  padded line  ") {
		t.Errorf("output should trim leading/trailing whitespace from unparsed lines:\n%s", out)
	}
}

func TestWriteReport_EmptyResult(t *testing.T) {
	result := &Result{}

	var buf bytes.Buffer
	if err := result.WriteReport(&buf, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "그룹 0개 (미분류 0줄)") {
		t.Errorf("output missing zero-count header:\n%s", out)
	}
	if !strings.Contains(out, "미분류 줄 (0개)") {
		t.Errorf("output missing zero-count unparsed header:\n%s", out)
	}
}

func TestWriteGroupTable_IncludesSequentialEventIDs(t *testing.T) {
	observed := time.Date(
		2026,
		time.June,
		1,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	result := &Result{
		Groups: []Group{
			{
				Facility: "ETHPORT",
				Mnemonic: "IF_DOWN_LINK_FAILURE",
				Iface:    "IPStorage1/6",
				Vsan:     "-",
				Severity: "5",
				Count:    1,
				First:    observed,
				Last:     observed,
			},
			{
				Facility: "ETHPORT",
				Mnemonic: "IF_SFP_WARNING",
				Iface:    "IPStorage1/6",
				Vsan:     "-",
				Severity: "4",
				Count:    1,
				First:    observed,
				Last:     observed,
			},
			{
				Facility: "ETHPORT",
				Mnemonic: "IF_SFP_ALARM",
				Iface:    "IPStorage1/6",
				Vsan:     "-",
				Severity: "3",
				Count:    1,
				First:    observed,
				Last:     observed,
			},
		},
	}

	var buf bytes.Buffer
	if err := result.WriteGroupTable(&buf); err != nil {
		t.Fatalf("WriteGroupTable returned an error: %v", err)
	}

	output := buf.String()

	expected := []string{
		"[E1 sev5] ETHPORT-IF_DOWN_LINK_FAILURE",
		"[E2 sev4] ETHPORT-IF_SFP_WARNING",
		"[E3 sev3] ETHPORT-IF_SFP_ALARM",
	}

	for _, value := range expected {
		if !strings.Contains(output, value) {
			t.Errorf(
				"output does not contain %q:\n%s",
				value,
				output,
			)
		}
	}
}
