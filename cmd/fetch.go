package cmd

import (
	"github.com/lnsp/ftp2p/pkg/fetcher"
	"github.com/lnsp/ftp2p/pkg/tracker"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [tracker] [destination]",
	Short: "Fetch file from peer-to-peer network",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		tpath, path := args[0], args[1]
		t, err := tracker.Open(tpath)
		if err != nil {
			return err
		}
		if err := fetcher.Fetch(path, t); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
