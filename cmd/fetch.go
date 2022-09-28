package cmd

import (
	"runtime"

	"github.com/lnsp/ftp2p/fetcher"
	"github.com/lnsp/ftp2p/tracker"

	"github.com/spf13/cobra"
)

var (
	fetchTracker  string
	fetchOutput   string
	fetchWorkers  int
	fetchInsecure bool
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch file from peer-to-peer network",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := tracker.Open(fetchTracker)
		if err != nil {
			return err
		}
		if err := fetcher.Fetch(fetchOutput, t, fetchWorkers, fetchInsecure); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	fetchCmd.Flags().StringVarP(&fetchTracker, "tracker", "t", "tracker", "Tracker file to use")
	fetchCmd.Flags().StringVarP(&fetchOutput, "output", "o", "output", "Output file path")
	fetchCmd.Flags().IntVarP(&fetchWorkers, "workers", "n", runtime.NumCPU(), "Number of workers used for pulling chunks")
	fetchCmd.Flags().BoolVarP(&fetchInsecure, "insecure", "x", false, "Allow insecure connection to Tavern")
	rootCmd.AddCommand(fetchCmd)
}
