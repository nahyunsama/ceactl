package logcompressor

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	reTimestamp = regexp.MustCompile(`^(\d{4}\s+\w+\s+\d+\s+\d{2}:\d{2}:\d{2})`)
	reMnemonic  = regexp.MustCompile(`%([A-Z0-9_-]+)-(\d)-([A-Z0-9_]+):`)
	reInterface = regexp.MustCompile(`Interface (\S+?)(?:,|\s+is\b)`)
	reVsan      = regexp.MustCompile(`(?i)VSAN\s*(\d+)`)
)

type groupKey struct {
	facility string
	mnemonic string
	iface    string
	vsan     string
}

type Group struct {
	Facility string
	Mnemonic string
	Iface    string
	Vsan     string
	Severity string
	Sample   string
	Count    int
	First    time.Time
	Last     time.Time
}

type Result struct {
	Groups   []Group
	Unparsed []string
}

func parseTS(line string) (time.Time, bool) {
	m := reTimestamp.FindStringSubmatch(line)
	if m == nil {
		return time.Time{}, false
	}
	normalized := strings.Join(strings.Fields(m[1]), " ")
	t, err := time.Parse("2006 Jan 2 15:04:05", normalized)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func parseMnemonic(line string) (facility, severity, mnemonic string, ok bool) {
	m := reMnemonic.FindStringSubmatch(line)
	if m == nil {
		return "", "", "", false
	}
	return m[1], m[2], m[3], true
}

func parseInterface(line string) string {
	m := reInterface.FindStringSubmatch(line)
	if m == nil {
		return "-"
	}
	return m[1]
}

func parseVsan(line string) string {
	m := reVsan.FindStringSubmatch(line)
	if m == nil {
		return "-"
	}
	return m[1]
}

func Analyze(r io.Reader, from, to time.Time) (*Result, error) {
	groups := make(map[groupKey]*Group)
	var order []groupKey
	var unparsed []string

	scanner := bufio.NewScanner(r)

	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		ts, tsOK := parseTS(line)
		if !tsOK {
			unparsed = append(unparsed, line)
			continue
		}
		if !from.IsZero() && ts.Before(from) {
			continue
		}
		if !to.IsZero() && ts.After(to) {
			continue
		}

		facility, severity, mnemonic, mOK := parseMnemonic(line)
		if !mOK {
			unparsed = append(unparsed, line)
			continue
		}

		iface := parseInterface(line)
		vsan := parseVsan(line)
		key := groupKey{facility: facility, mnemonic: mnemonic, iface: iface, vsan: vsan}

		g, exists := groups[key]
		if !exists {
			g = &Group{
				Facility: facility,
				Mnemonic: mnemonic,
				Iface:    iface,
				Vsan:     vsan,
				Severity: severity,
				Sample:   strings.TrimSpace(line),
				Count:    0,
				First:    ts,
				Last:     ts,
			}
			groups[key] = g
			order = append(order, key)
		}
		g.Count++
		if ts.Before(g.First) {
			g.First = ts
		}
		if ts.After(g.Last) {
			g.Last = ts
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Slice(order, func(i, j int) bool {
		return groups[order[i]].First.Before(groups[order[j]].First)
	})

	result := &Result{Unparsed: unparsed}
	for _, key := range order {
		result.Groups = append(result.Groups, *groups[key])
	}
	return result, nil
}

func (r *Result) WriteReport(w io.Writer, maxUnparsed int) error {
	var err error
	write := func(format string, a ...any) {
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(w, format, a...)
	}

	write("=== 압축 결과: 그룹 %d개 (미분류 %d줄) ===\n\n", len(r.Groups), len(r.Unparsed))
	for _, g := range r.Groups {
		var span string
		if !g.First.Equal(g.Last) {
			span = fmt.Sprintf("%s ~ %s", g.First.Format("15:04:05"), g.Last.Format("15:04:05"))
		} else {
			span = g.First.Format("15:04:05")
		}
		write("[sev%s] %s-%s  iface=%s vsan=%s  %d회  (%s)\n",
			g.Severity, g.Facility, g.Mnemonic, g.Iface, g.Vsan, g.Count, span)
		write("    예시 원문: %s\n", g.Sample)
	}

	write("\n=== 미분류 줄 (%d개) ===\n", len(r.Unparsed))
	limit := maxUnparsed
	if len(r.Unparsed) < limit {
		limit = len(r.Unparsed)
	}
	for _, l := range r.Unparsed[:limit] {
		write("  %s\n", strings.TrimSpace(l))
	}
	if len(r.Unparsed) > limit {
		write("  ... 외 %d줄\n", len(r.Unparsed)-limit)
	}
	return err
}
