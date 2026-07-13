package llmanalysis

import (
	"bytes"
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/nahyunsama/ceactl/internal/mds/logcompressor"
)

//go:embed prompts/system_prompt.txt
var SystemPrompt string

var reRepeated = regexp.MustCompile(
	`(?i)\blast message repeated\s+(\d+)\s+times?\b`,
)

type PromptInput struct {
	Device      string
	FilterStart time.Time
	FilterEnd   time.Time
	Result      *logcompressor.Result
}

type UnparsedSummary struct {
	RepeatNoticeLines           int
	UnassignedRepeatOccurrences int
	OtherLines                  int
}

func BuildUserPrompt(input PromptInput) (string, error) {
	if input.Result == nil {
		return "", fmt.Errorf("BuildUserPrompt: result is nil")
	}

	var buf bytes.Buffer

	buf.WriteString(
		"Analyze the following grouped Cisco device log events.\n" +
			"Follow the system instructions. \n\n",
	)

	writeMetadata(&buf, input)
	writeEvents(&buf, input.Result.Groups)
	writeUnparsed(&buf, summarizeUnparsed(input.Result.Unparsed))

	buf.WriteString(
		"<notes>\n" +
			"Each event row was grouped mechanically by " +
			"(facility, mnemonic, interface, vsan).\n" +
			"observed_count contains only parsed messages assigned to that row.\n" +
			"Unassigned repeat occurrences are not included in observed_count.\n" +
			"Unassigned repeat occurrences do not belong to a specific event row.\n" +
			"</notes>\n",
	)

	return buf.String(), nil
}

func writeMetadata(buf *bytes.Buffer, input PromptInput) {
	start, end, ok := eventRange(input.Result.Groups)

	buf.WriteString("<metadata>\n")
	fmt.Fprintf(buf, "device: %s\n", quotedOrNull(input.Device))
	buf.WriteString(`source_command: "show logging logfile"` + "\n")
	fmt.Fprintf(buf, "filter_start: %s\n", formattedTimeOrNull(input.FilterStart))
	fmt.Fprintf(buf, "filter_end: %s\n", formattedTimeOrNull(input.FilterEnd))

	if ok {
		fmt.Fprintf(buf, "event_time_min: %s\n", formattedTimeOrNull(start))
		fmt.Fprintf(buf, "event_time_max: %s\n", formattedTimeOrNull(end))
	} else {
		buf.WriteString("event_time_min: null\n")
		buf.WriteString("event_time_max: null\n")
	}

	buf.WriteString(
		`timestamp_basis: "device log time; timezone not provided"` + "\n",
	)
	buf.WriteString("</metadata>\n\n")
}

func writeEvents(buf *bytes.Buffer, groups []logcompressor.Group) {
	fmt.Fprintf(buf, "<events count=%q>\n", strconv.Itoa(len(groups)))

	for i, group := range groups {
		fmt.Fprintf(
			buf,
			"event id=%d severity=%s facility=%s mnemonic=%s "+
				"interface=%s vsan=%s observed_count=%d first=%s last=%s\n",
			i+1,
			quotedOrNull(group.Severity),
			quotedOrNull(group.Facility),
			quotedOrNull(group.Mnemonic),
			nullableParsedValue(group.Iface),
			nullableParsedValue(group.Vsan),
			group.Count,
			formattedTimeOrNull(group.First),
			formattedTimeOrNull(group.Last),
		)
	}

	buf.WriteString("</events>\n\n")
}

func writeUnparsed(buf *bytes.Buffer, summary UnparsedSummary) {
	buf.WriteString("<unparsed>\n")
	fmt.Fprintf(
		buf,
		"repeat_notice_lines: %d\n",
		summary.RepeatNoticeLines,
	)
	fmt.Fprintf(
		buf,
		"unassigned_repeat_occurrences: %d\n",
		summary.UnassignedRepeatOccurrences,
	)
	fmt.Fprintf(
		buf,
		"other_unparsed_lines: %d\n",
		summary.OtherLines,
	)
	buf.WriteString("</unparsed>\n\n")
}

func summarizeUnparsed(lines []string) UnparsedSummary {
	var summary UnparsedSummary

	for _, line := range lines {
		match := reRepeated.FindStringSubmatch(line)
		if match == nil {
			summary.OtherLines++
			continue
		}

		count, err := strconv.Atoi(match[1])
		if err != nil {
			summary.OtherLines++
			continue
		}

		summary.RepeatNoticeLines++
		summary.UnassignedRepeatOccurrences += count
	}

	return summary
}

func eventRange(
	groups []logcompressor.Group,
) (start time.Time, end time.Time, ok bool) {
	if len(groups) == 0 {
		return time.Time{}, time.Time{}, false
	}

	start = groups[0].First
	end = groups[0].Last

	for _, group := range groups[1:] {
		if group.First.Before(start) {
			start = group.First
		}
		if group.Last.After(end) {
			end = group.Last
		}
	}
	return start, end, true
}

func quotedOrNull(value string) string {
	if value == "" {
		return "null"
	}

	return strconv.Quote(value)
}

func nullableParsedValue(value string) string {
	if value == "" || value == "-" {
		return "null"
	}

	return strconv.Quote(value)
}

func formattedTimeOrNull(value time.Time) string {
	if value.IsZero() {
		return "null"
	}

	return strconv.Quote(value.Format("2006-01-02 15:04:05"))
}
