package cmd

import (
	"fmt"
	"strings"

	"github.com/s1na/hq/internal/api"
	"github.com/s1na/hq/internal/display"
	"github.com/spf13/cobra"
)

var (
	runsSim    string
	runsClient string
	runsLimit  int
)

var runsCmd = &cobra.Command{
	Use:   "runs",
	Short: "List recent test runs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		entries, err := client.FetchListing(runsSim, runsClient, 0)
		if err != nil {
			return fmt.Errorf("fetching listing: %w", err)
		}

		api.SortByTime(entries)

		if runsLimit > 0 && len(entries) > runsLimit {
			entries = entries[:runsLimit]
		}

		t := display.NewTable([]string{"Name", "Clients", "Tests", "Pass", "Fail", "When", "File"})
		for _, e := range entries {
			t.Append([]string{
				e.Name,
				strings.Join(e.Clients, ","),
				fmt.Sprintf("%d", e.NTests),
				fmt.Sprintf("%d", e.Passes),
				fmt.Sprintf("%d", e.Fails),
				api.FormatTime(e.Start),
				e.FileName,
			})
		}
		t.Render()
		return nil
	},
}

func init() {
	runsCmd.Flags().StringVar(&runsSim, "sim", "", "Filter by simulator name (substring match)")
	runsCmd.Flags().StringVar(&runsClient, "client", "", "Filter by client name")
	runsCmd.Flags().IntVar(&runsLimit, "limit", 20, "Maximum number of runs to show")
	rootCmd.AddCommand(runsCmd)
}
