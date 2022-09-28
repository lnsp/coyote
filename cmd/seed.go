package cmd

import (
	"fmt"

	"github.com/lnsp/ftp2p/security"
	"github.com/lnsp/ftp2p/seeder"
	"github.com/lnsp/ftp2p/tracker"

	"github.com/spf13/cobra"
)

var (
	seedAddr          string
	seedFile          string
	seedTracker       string
	seedAllowInsecure bool
	seedPoolsize      int
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed file for other to download",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		tlsConfig, err := security.NewBasicTLS()
		if err != nil {
			return fmt.Errorf("generate tls: %w", err)
		}
		track, err := tracker.Open(seedTracker)
		if err != nil {
			return err
		}
		seed := seeder.New(seedAddr, seedAllowInsecure, seedPoolsize, tlsConfig)
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
	seedCmd.Flags().IntVarP(&seedPoolsize, "poolsize", "p", 64, "Size of the concurrent connection pool")
	seedCmd.Flags().BoolVarP(&seedAllowInsecure, "insecure", "x", false, "Allow insecure connection to Tavern")
	rootCmd.AddCommand(seedCmd)
}
