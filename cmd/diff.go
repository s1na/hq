package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/s1na/hq/internal/api"
	"github.com/s1na/hq/internal/display"
	"github.com/spf13/cobra"
)

var (
	diffTest   string
	diffClient string
)

var diffCmd = &cobra.Command{
	Use:   "diff <run-file>",
	Short: "Show colorized diff for failing tests",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.FetchResult(args[0])
		if err != nil {
			return fmt.Errorf("fetching result: %w", err)
		}

		if result.TestDetailsLog == "" {
			return fmt.Errorf("no test details log available for this run")
		}

		matched := 0
		for _, tc := range result.TestCases {
			if tc.SummaryResult.Pass {
				continue
			}

			clientName := api.ExtractClient(tc.Name)
			if diffClient != "" && !strings.Contains(strings.ToLower(clientName), strings.ToLower(diffClient)) {
				continue
			}

			if diffTest != "" {
				m, _ := path.Match(diffTest, tc.Name)
				if !m {
					// Try matching without client suffix
					parts := strings.TrimSpace(tc.Name)
					idx := strings.LastIndex(parts, " (")
					if idx > 0 {
						m, _ = path.Match(diffTest, parts[:idx])
					}
				}
				if !m {
					continue
				}
			}

			begin := tc.SummaryResult.Log.Begin
			end := tc.SummaryResult.Log.End
			if begin == 0 && end == 0 {
				continue
			}

			log, err := client.FetchTestLog(result.TestDetailsLog, begin, end)
			if err != nil {
				fmt.Printf("Error fetching log for %s: %v\n", tc.Name, err)
				continue
			}

			matched++
			display.Bold.Printf("=== %s ===\n", tc.Name)
			display.ColorizeDiff(log, noColor)
			fmt.Println()
		}

		if matched == 0 {
			fmt.Println("No matching failing tests with log data found.")
		}
		return nil
	},
}

func init() {
	diffCmd.Flags().StringVar(&diffTest, "test", "", "Filter by test name (glob pattern)")
	diffCmd.Flags().StringVar(&diffClient, "client", "", "Filter by client name")
	rootCmd.AddCommand(diffCmd)
}
