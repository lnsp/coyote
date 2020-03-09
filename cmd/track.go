package cmd

import (
	"github.com/lnsp/ftp2p/pkg/tracker"

	"github.com/spf13/cobra"
)

var (
	trackInput   string
	trackTracker string
	trackTavern  string
)

var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Create a file tracker with a specified host tavern",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := tracker.Track(trackTavern, trackInput)
		if err != nil {
			return err
		}
		return tracker.Save(t, trackTracker)
	},
}

func init() {
	trackCmd.Flags().StringVarP(&trackInput, "input", "i", "input", "Input file to track")
	trackCmd.Flags().StringVarP(&trackTavern, "tavern", "s", "localhost:6443", "Tavern host which serves the peer sessions")
	trackCmd.Flags().StringVarP(&trackTracker, "tracker", "t", "tracker", "Tracker file")
	rootCmd.AddCommand(trackCmd)
}
