package cmd

import (
	"github.com/lnsp/ftp2p/pkg/fetcher"
	"github.com/lnsp/ftp2p/pkg/tracker"
	"runtime"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [tracker] [destination]",
	Short: "Fetch file from peer-to-peer network",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		tpath, path := args[0], args[1]
		n, err := cmd.Flags().GetInt("workers")
		if err != nil {
			return err
		}
		t, err := tracker.Open(tpath)
		if err != nil {
			return err
		}
		if err := fetcher.Fetch(path, t, n); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	fetchCmd.Flags().IntP("workers", "n", runtime.NumCPU(), "Number of workers used for pulling chunks")
	rootCmd.AddCommand(fetchCmd)
}
