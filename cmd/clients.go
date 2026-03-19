package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/s1na/hq/internal/api"
	"github.com/s1na/hq/internal/display"
	"github.com/spf13/cobra"
)

var clientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "List known clients and their aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		// Fetch recent runs to discover active clients.
		entries, err := client.FetchListing("", "", 0)
		if err != nil {
			return fmt.Errorf("fetching listing: %w", err)
		}

		// Collect unique client names.
		seen := make(map[string]bool)
		for _, e := range entries {
			for _, c := range e.Clients {
				seen[c] = true
			}
		}

		// Build reverse map: canonical prefix -> list of aliases.
		aliases := api.ClientAliases()
		reverseAliases := make(map[string][]string)
		for alias, canon := range aliases {
			if alias != canon {
				reverseAliases[canon] = append(reverseAliases[canon], alias)
			}
		}
		for _, v := range reverseAliases {
			sort.Strings(v)
		}

		// Sort client names.
		var clients []string
		for c := range seen {
			clients = append(clients, c)
		}
		sort.Strings(clients)

		t := display.NewTable([]string{"Client", "Alias"})
		for _, c := range clients {
			// Find alias by matching prefix.
			prefix := strings.TrimSuffix(c, "_default")
			prefix = strings.Split(prefix, "_")[0]
			aliasStr := ""
			if shortNames, ok := reverseAliases[prefix]; ok {
				aliasStr = strings.Join(shortNames, ", ")
			}
			t.Append([]string{c, aliasStr})
		}
		t.Render()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clientsCmd)
}
