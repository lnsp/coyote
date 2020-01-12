package cmd

import "github.com/spf13/cobra"

var rootCmd = cobra.Command{
	Use:   "ftp2p",
	Short: "Secure peer-to-peer file-transfer.",
}

func Execute() error {
	return rootCmd.Execute()
}
