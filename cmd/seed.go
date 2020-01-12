package cmd

import (
	"github.com/lnsp/ftp2p/pkg/seeder"
	"github.com/lnsp/ftp2p/pkg/tracker"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed [addr] [file] [tracker]",
	Short: "Seed file for other to download",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, path, tpath := args[0], args[1], args[2]
		track, err := tracker.Open(tpath)
		if err != nil {
			return err
		}
		seed := seeder.New(addr, 64)
		if err := seed.Seed(path, track); err != nil {
			return err
		}
		if err := seed.ListenAndServe(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
