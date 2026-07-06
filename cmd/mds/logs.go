package mds

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	appconfig "github.com/nahyunsama/ceactl/internal/config"
	"github.com/nahyunsama/ceactl/internal/mds/commands"
	"github.com/nahyunsama/ceactl/internal/mds/config"
	"github.com/nahyunsama/ceactl/internal/mds/logcompressor"
	"github.com/spf13/cobra"
)

func LogsCommand(opts *commandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "MDS log related commands",
	}

	cmd.AddCommand(LogsAnalyzeCommand(opts))

	return cmd
}

func LogsAnalyzeCommand(opts *commandOptions) *cobra.Command {
	var fromStr, toStr, file string

	c := &cobra.Command{
		Use:   "analyze",
		Short: "Group and summarize MDS log lines",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgFile, err := appconfig.LoadFile(opts.configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %v", err)
			}
			if !cfgFile.LLMAnalysis.Enabled {
				return fmt.Errorf("log analysis is disabled; set `llm_analysis.enabled: true` in %s to use `mds logs analyze`", opts.configPath)
			}

			from, err := parseDayStart(fromStr)
			if err != nil {
				return fmt.Errorf("invalid --from: %v", err)
			}
			to, err := parseDayEnd(toStr)
			if err != nil {
				return fmt.Errorf("invalid --to: %v", err)
			}

			var reader io.Reader

			if file != "" {
				f, err := os.Open(file)
				if err != nil {
					return fmt.Errorf("failed to open log file: %v", err)
				}
				defer f.Close()
				reader = f
			} else {
				cfg, err := config.LoadConfig(opts.configPath, opts.deviceName, opts.verbose)
				if err != nil {
					return fmt.Errorf("failed to load config: %v", err)
				}

				body, err := commands.GetLoggingLogfile(cmd.Context(), cfg)
				if err != nil {
					return fmt.Errorf("failed to get logging logfile: %v", err)
				}
				reader = strings.NewReader(body)
			}

			result, err := logcompressor.Analyze(reader, from, to)
			if err != nil {
				return fmt.Errorf("failed to analyze log: %v", err)
			}

			return result.WriteReport(os.Stdout, 10)
		},
	}

	c.Flags().StringVar(&fromStr, "from", "", "start date filter, YYYYMMDD (inclusive)")
	c.Flags().StringVar(&toStr, "to", "", "end date filter, YYYYMMDD (inclusive)")
	c.Flags().StringVar(&file, "file", "", "path to a local log file (skips device fetch)")

	return c
}

func parseDayStart(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse("20060102", s)
}

func parseDayEnd(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	day, err := time.Parse("20060102", s)
	if err != nil {
		return time.Time{}, err
	}
	return day.AddDate(0, 0, 1).Add(-time.Second), nil
}
