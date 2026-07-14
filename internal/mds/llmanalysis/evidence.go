package llmanalysis

import (
	"regexp"
	"strconv"
)

var (
	eventRangeReference = regexp.MustCompile(
		`(?i)\bE\d+\s*(?:-|through|to)\s*E?\d+\b`,
	)
	eventReference = regexp.MustCompile(`\bE([1-9]\d*)\b`)
)

// ReferencedEventIDs returns valid, individually cited Event IDs in their
// first-seen order. Range references are excluded because they do not prove
// that each event in the range directly supports the surrounding statement.
func ReferencedEventIDs(reply string, groupCount int) []int {
	if groupCount <= 0 {
		return nil
	}

	withoutRanges := eventRangeReference.ReplaceAllString(reply, " ")
	matches := eventReference.FindAllStringSubmatch(withoutRanges, -1)

	seen := make(map[int]struct{}, len(matches))
	ids := make([]int, 0, len(matches))

	for _, match := range matches {
		id, err := strconv.Atoi(match[1])
		if err != nil || id > groupCount {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}

		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	return ids
}
