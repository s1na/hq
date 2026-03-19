package cmd

import (
	"fmt"
	"strings"

	"github.com/s1na/hq/internal/api"
	"github.com/s1na/hq/internal/display"
	"github.com/spf13/cobra"
)

var (
	statsSim    string
	statsClient string
	statsLast   int
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show pass/fail rates across runs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if statsSim == "" {
			return fmt.Errorf("--sim flag is required")
		}

		client, err := newClient()
		if err != nil {
			return err
		}

		entries, err := client.FetchAllListing(statsSim)
		if err != nil {
			return fmt.Errorf("fetching listing: %w", err)
		}

		api.SortByTime(entries)

		if statsLast > 0 && len(entries) > statsLast {
			entries = entries[:statsLast]
		}

		if statsClient != "" {
			// Show per-run stats for a specific client
			t := display.NewTable([]string{"Run", "Tests", "Pass", "Fail", "Rate", "When"})
			for _, e := range entries {
				if !containsClientCI(e.Clients, statsClient) {
					continue
				}
				// For per-client stats, we need to fetch the result
				result, err := client.FetchResult(e.FileName)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: skipping %s: %v\n", e.FileName, err)
					continue
				}

				passes, fails := countForClient(result, statsClient)
				total := passes + fails
				rate := "N/A"
				if total > 0 {
					rate = fmt.Sprintf("%.1f%%", float64(passes)/float64(total)*100)
				}
				t.Append([]string{
					e.FileName,
					fmt.Sprintf("%d", total),
					fmt.Sprintf("%d", passes),
					fmt.Sprintf("%d", fails),
					rate,
					api.FormatTime(e.Start),
				})
			}
			t.Render()
		} else {
			// Show aggregate stats per client across runs
			type clientStats struct {
				passes int
				fails  int
				runs   int
			}
			stats := make(map[string]*clientStats)

			for _, e := range entries {
				result, err := client.FetchResult(e.FileName)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: skipping %s: %v\n", e.FileName, err)
					continue
				}

				clientCounts := make(map[string][2]int) // [passes, fails]
				for _, tc := range result.TestCases {
					cl := api.ExtractClient(tc.Name)
					if cl == "" {
						continue
					}
					counts := clientCounts[cl]
					if tc.SummaryResult.Pass {
						counts[0]++
					} else {
						counts[1]++
					}
					clientCounts[cl] = counts
				}

				for cl, counts := range clientCounts {
					s, ok := stats[cl]
					if !ok {
						s = &clientStats{}
						stats[cl] = s
					}
					s.passes += counts[0]
					s.fails += counts[1]
					s.runs++
				}
			}

			t := display.NewTable([]string{"Client", "Runs", "Total Tests", "Pass", "Fail", "Rate"})
			for cl, s := range stats {
				total := s.passes + s.fails
				rate := "N/A"
				if total > 0 {
					rate = fmt.Sprintf("%.1f%%", float64(s.passes)/float64(total)*100)
				}
				t.Append([]string{
					cl,
					fmt.Sprintf("%d", s.runs),
					fmt.Sprintf("%d", total),
					fmt.Sprintf("%d", s.passes),
					fmt.Sprintf("%d", s.fails),
					rate,
				})
			}
			t.Render()
		}

		return nil
	},
}

func countForClient(result *api.TestSuiteResult, clientName string) (passes, fails int) {
	for _, tc := range result.TestCases {
		cl := api.ExtractClient(tc.Name)
		if !strings.Contains(strings.ToLower(cl), strings.ToLower(clientName)) {
			continue
		}
		if tc.SummaryResult.Pass {
			passes++
		} else {
			fails++
		}
	}
	return
}

func containsClientCI(clients []string, target string) bool {
	target = strings.ToLower(target)
	for _, c := range clients {
		if strings.Contains(strings.ToLower(c), target) {
			return true
		}
	}
	return false
}

func init() {
	statsCmd.Flags().StringVar(&statsSim, "sim", "", "Simulator name (required)")
	statsCmd.Flags().StringVar(&statsClient, "client", "", "Filter by client name")
	statsCmd.Flags().IntVar(&statsLast, "last", 10, "Number of recent runs to analyze")
	rootCmd.AddCommand(statsCmd)
}
