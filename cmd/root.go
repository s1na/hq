package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/s1na/hq/internal/api"
	"github.com/s1na/hq/internal/cache"
	"github.com/spf13/cobra"
)

var (
	baseURL  string
	suite    string
	noCache  bool
	noColor  bool
	cacheDir string
)

var rootCmd = &cobra.Command{
	Use:   "hq",
	Short: "hq (hive query) - query and investigate Ethereum Hive test results",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if noColor {
			color.NoColor = true
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "https://hive.ethpandaops.io", "Hive server base URL")
	rootCmd.PersistentFlags().StringVar(&suite, "suite", "generic", "Test suite name")
	rootCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false, "Bypass cache reads")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", "", "Cache directory (default ~/.cache/hq)")
}

func newClient() (*api.Client, error) {
	c, err := cache.New(cacheDir, !noCache)
	if err != nil {
		return nil, fmt.Errorf("initializing cache: %w", err)
	}
	return api.NewClient(baseURL, suite, c), nil
}
