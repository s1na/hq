package cmd

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/s1na/hq/internal/api"
	"github.com/s1na/hq/internal/display"
	"github.com/spf13/cobra"
)

var (
	testsSim    string
	testsClient string
	testsFilter string
)

var testsCmd = &cobra.Command{
	Use:   "tests [run-file]",
	Short: "List test cases with pass/fail status",
	Long: `List test cases with their pass/fail status. If no run-file is given,
uses the most recent run matching the --sim and --client filters.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		testsClient = api.ResolveClientAlias(testsClient)

		client, err := newClient()
		if err != nil {
			return err
		}

		var fileName string
		if len(args) > 0 {
			fileName = args[0]
		} else {
			entries, err := client.FetchListing(testsSim, testsClient, 0)
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

		type testEntry struct {
			name   string
			client string
			pass   bool
		}

		var tests []testEntry
		for _, tc := range result.TestCases {
			clientName := api.ExtractClient(tc.Name)
			if testsClient != "" && !strings.Contains(strings.ToLower(clientName), strings.ToLower(testsClient)) {
				continue
			}
			if testsFilter != "" {
				m, _ := path.Match(testsFilter, tc.Name)
				if !m {
					parts := strings.TrimSpace(tc.Name)
					idx := strings.LastIndex(parts, " (")
					if idx > 0 {
						m, _ = path.Match(testsFilter, parts[:idx])
					}
				}
				if !m {
					continue
				}
			}
			tests = append(tests, testEntry{
				name:   tc.Name,
				client: clientName,
				pass:   tc.SummaryResult.Pass,
			})
		}

		if len(tests) == 0 {
			fmt.Println("No matching tests found.")
			return nil
		}

		sort.Slice(tests, func(i, j int) bool {
			return tests[i].name < tests[j].name
		})

		passes := 0
		t := display.NewTable([]string{"Test", "Status"})
		for _, tc := range tests {
			status := display.PassFail(tc.pass)
			t.Append([]string{tc.name, status})
			if tc.pass {
				passes++
			}
		}
		t.Render()

		fmt.Printf("\n%s (%d/%d passing)\n", display.PassFailCount(passes, len(tests)), passes, len(tests))
		return nil
	},
}

func init() {
	testsCmd.Flags().StringVar(&testsSim, "sim", "", "Filter runs by simulator name")
	testsCmd.Flags().StringVar(&testsClient, "client", "", "Filter by client name")
	testsCmd.Flags().StringVar(&testsFilter, "test", "", "Filter by test name (glob pattern)")
	rootCmd.AddCommand(testsCmd)
}
