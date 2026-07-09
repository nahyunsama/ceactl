package llmanalysis

import (
	"bytes"
	_ "embed"
	"fmt"
	"regexp"
	"strconv"

	"github.com/nahyunsama/ceactl/internal/mds/logcompressor"
)

//go:embed prompts/system_prompt.txt
var SystemPrompt string

var reRepeated = regexp.MustCompile(`repeated (\d+) times`)

func BuildUserPrompt(device string, result *logcompressor.Result) (string, error) {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "Below is a compressed event-group summary from a Cisco MDS device (%s), covering %s.\n", device, formatRange(result.Groups))
	fmt.Fprintf(&buf, "Grouping was done mechanically by (facility, mnemonic, interface, vsan) — no interpretation has been applied.\n\n")
	fmt.Fprintf(&buf, "Format: [severity] facility-MNEMONIC  iface=X vsan=Y  COUNT회  (first ~ last)\n\n")

	if err := result.WriteGroupTable(&buf); err != nil {
		return "", fmt.Errorf("BuildUserPrompt: %w", err)
	}

	if note := unparsedNote(result.Unparsed); note != "" {
		fmt.Fprintf(&buf, "\n%s\n", note)
	}

	return buf.String(), nil
}

func formatRange(groups []logcompressor.Group) string {
	if len(groups) == 0 {
		return "an unknown time range"
	}

	first, last := groups[0].First, groups[0].Last
	for _, g := range groups[1:] {
		if g.First.Before(first) {
			first = g.First
		}
		if g.Last.After(last) {
			last = g.Last
		}
	}

	if first.Format("2006-01-02") == last.Format("2006-01-02") {
		return fmt.Sprintf("%s, %s–%s", first.Format("2006 Jan 2"), first.Format("15:04:05"), last.Format("15:04:05"))
	}
	return fmt.Sprintf("%s–%s", first.Format("2006 Jan 2, 15:04:05"), last.Format("2006 Jan 2, 15:04:05"))
}

func unparsedNote(unparsed []string) string {
	if len(unparsed) == 0 {
		return ""
	}

	total := 0
	for _, line := range unparsed {
		if m := reRepeated.FindStringSubmatch(line); m != nil {
			n, _ := strconv.Atoi(m[1])
			total += n
		}
	}

	return fmt.Sprintf(
		"Note: %d lines in the original log were device-side \"last message repeated N times\" compressions "+
			"(totaling ~%d additional occurrences) and are NOT included in the counts above. "+
			"Treat the counts as a floor, not an exact total.",
		len(unparsed), total,
	)
}
