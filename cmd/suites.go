package cmd

import (
	"fmt"

	"github.com/s1na/hq/internal/display"
	"github.com/spf13/cobra"
)

var suitesCmd = &cobra.Command{
	Use:   "suites",
	Short: "List available test suites",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		discoveries, err := client.FetchDiscovery()
		if err != nil {
			return fmt.Errorf("fetching discovery: %w", err)
		}

		t := display.NewTable([]string{"Name", "Address"})
		for _, d := range discoveries {
			t.Append([]string{d.Name, d.Address})
		}
		t.Render()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(suitesCmd)
}
