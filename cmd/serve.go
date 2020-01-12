package cmd

import (
	"github.com/lnsp/ftp2p/pkg/tavern"

	"github.com/spf13/cobra"
)

var announcementTimeout int64

var serveCmd = &cobra.Command{
	Use:   "serve [addr]",
	Short: "Serve as a tavern host and track peers",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		server := tavern.New(announcementTimeout)
		return server.ListenAndServe(args[0])
	},
}

func init() {
	serveCmd.Flags().Int64VarP(&announcementTimeout, "timeout", "t", 300, "Announcement timeout in seconds")
	rootCmd.AddCommand(serveCmd)
}
