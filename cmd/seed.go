package cmd

import (
	"github.com/lnsp/ftp2p/pkg/seeder"
	"github.com/lnsp/ftp2p/pkg/tracker"

	"github.com/spf13/cobra"
)

var (
	seedAddr    string
	seedFile    string
	seedTracker string
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed file for other to download",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		track, err := tracker.Open(seedTracker)
		if err != nil {
			return err
		}
		seed := seeder.New(seedAddr, 64)
		if err := seed.Seed(seedFile, track); err != nil {
			return err
		}
		if err := seed.ListenAndServe(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	seedCmd.Flags().StringVarP(&seedAddr, "listen", "l", "localhost:6444", "Address to listen on")
	seedCmd.Flags().StringVarP(&seedFile, "input", "i", "input", "Tracked input file to seed")
	seedCmd.Flags().StringVarP(&seedTracker, "tracker", "t", "tracker", "Tracker to use")
	rootCmd.AddCommand(seedCmd)
}
