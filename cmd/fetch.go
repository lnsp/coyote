package cmd

import (
	"runtime"

	"github.com/lnsp/ftp2p/fetcher"
	"github.com/lnsp/ftp2p/tracker"

	"github.com/spf13/cobra"
)

var (
	fetchOutput   string
	fetchWorkers  int
	fetchInsecure bool
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [tracker]",
	Short: "Fetch file from peer-to-peer network",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tracker, err := tracker.Open(args[0])
		if err != nil {
			return err
		}
		fetcher := fetcher.New(fetchWorkers, fetchInsecure)
		if err := fetcher.Fetch(fetchOutput, tracker); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	fetchCmd.Flags().StringVarP(&fetchOutput, "output", "o", "output", "Output file path")
	fetchCmd.Flags().IntVarP(&fetchWorkers, "workers", "n", runtime.NumCPU(), "Number of workers used for pulling chunks")
	fetchCmd.Flags().BoolVarP(&fetchInsecure, "insecure", "x", false, "Allow insecure connection to Tavern")
	rootCmd.AddCommand(fetchCmd)
}
