package cmd

import (
	"github.com/lnsp/ftp2p/pkg/tracker"

	"github.com/spf13/cobra"
)

var trackCmd = &cobra.Command{
	Use:   "track [file] [tavern]",
	Short: "Create a file tracker with a specified host tavern",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, tavern := args[0], args[1]
		t, err := tracker.Track(tavern, path)
		if err != nil {
			return err
		}
		tpath := path + ".tracker"
		return tracker.Save(t, tpath)
	},
}

func init() {
	rootCmd.AddCommand(trackCmd)
}
