package cmd

import (
	"github.com/lnsp/ftp2p/pkg/tavern"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve [addr]",
	Short: "Serve as a tavern host and track peers",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		server := tavern.New()
		return server.ListenAndServe(args[0])
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
