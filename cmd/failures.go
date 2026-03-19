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
	failuresSim    string
	failuresClient string
	failuresTest   string
)

var failuresCmd = &cobra.Command{
	Use:   "failures [run-file]",
	Short: "Show failing tests from a run",
	Long: `Show failing tests. If no run-file is given, uses the most recent run
matching the --sim and --client filters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		var fileName string
		if len(args) > 0 {
			fileName = args[0]
		} else {
			// Find most recent matching run
			entries, err := client.FetchListing(failuresSim, failuresClient, 0)
			if err != nil {
				return fmt.Errorf("fetching listing: %w", err)
			}
			if len(entries) == 0 {
				return fmt.Errorf("no runs found matching filters")
			}
			api.SortByTime(entries)
			fileName = entries[0].FileName
			fmt.Printf("Using most recent run: %s\n\n", fileName)
		}

		result, err := client.FetchResult(fileName)
		if err != nil {
			return fmt.Errorf("fetching result: %w", err)
		}

		// Collect failures grouped by client
		type failure struct {
			id       string
			name     string
			client   string
			details  string
		}

		var failures []failure
		for id, tc := range result.TestCases {
			if tc.SummaryResult.Pass {
				continue
			}
			clientName := api.ExtractClient(tc.Name)
			if failuresClient != "" && !strings.Contains(strings.ToLower(clientName), strings.ToLower(failuresClient)) {
				continue
			}
			if failuresTest != "" {
				matched, _ := path.Match(failuresTest, tc.Name)
				if !matched {
					// Also try matching just the method/test part without the client suffix
					parts := strings.TrimSpace(tc.Name)
					idx := strings.LastIndex(parts, " (")
					if idx > 0 {
						matched, _ = path.Match(failuresTest, parts[:idx])
					}
				}
				if !matched {
					continue
				}
			}
			failures = append(failures, failure{
				id:      id,
				name:    tc.Name,
				client:  clientName,
				details: tc.SummaryResult.Details,
			})
		}

		if len(failures) == 0 {
			fmt.Println("No failures found.")
			return nil
		}

		fmt.Printf("Suite: %s\n", result.Name)
		fmt.Printf("Failures: %d\n\n", len(failures))

		t := display.NewTable([]string{"Client", "Test", "Details"})
		for _, f := range failures {
			details := f.details
			if len(details) > 80 {
				details = details[:77] + "..."
			}
			t.Append([]string{f.client, f.name, details})
		}
		t.Render()
		return nil
	},
}

func init() {
	failuresCmd.Flags().StringVar(&failuresSim, "sim", "", "Filter runs by simulator name")
	failuresCmd.Flags().StringVar(&failuresClient, "client", "", "Filter by client name")
	failuresCmd.Flags().StringVar(&failuresTest, "test", "", "Filter by test name (glob pattern)")
	rootCmd.AddCommand(failuresCmd)
}
