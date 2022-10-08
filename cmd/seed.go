package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/lnsp/coyote/security"
	"github.com/lnsp/coyote/seeder"
	"github.com/lnsp/coyote/tracker"
	"golang.org/x/net/context"

	"github.com/spf13/cobra"
)

var (
	seedAddr          string
	seedAllowInsecure bool
	seedPoolsize      int
)

var seedCmd = &cobra.Command{
	Use:   "seed [file] [(tracker)]",
	Short: "Seed file for other to download",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create sigcancel context
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt)
		// Create context waiting for sigs to be fired
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() { <-sigs; cancel() }()
		// Bootstrap TLS configuration
		tlsConfig, err := security.NewBasicTLS()
		if err != nil {
			return fmt.Errorf("generate tls: %w", err)
		}
		log.Println("Generated new ephemeral TLS certificate")
		// Use file and tracker path from args
		path := args[0]
		trackerpath := path + ".coyote"
		if len(args) > 1 {
			trackerpath = args[1]
		}
		// Open given tracker
		tracker, err := tracker.Open(trackerpath)
		if err != nil {
			return err
		}
		// Setup seeding instance
		log.Println("Setting up seeding instance listening on", seedAddr)
		seed := seeder.New(seedAddr, seedAllowInsecure, seedPoolsize, tlsConfig)
		if err := seed.Seed(ctx, path, tracker); err != nil {
			return err
		}
		log.Printf("Begin seeding %s to fetchers", hex.EncodeToString(tracker.Hash))
		if err := seed.ListenAndServe(ctx); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	seedCmd.Flags().StringVarP(&seedAddr, "listen", "l", "localhost:6444", "Address to listen on")
	seedCmd.Flags().IntVarP(&seedPoolsize, "poolsize", "p", 64, "Size of the concurrent connection pool")
	seedCmd.Flags().BoolVarP(&seedAllowInsecure, "insecure", "x", false, "Allow insecure connection to Tavern")
	rootCmd.AddCommand(seedCmd)
}
