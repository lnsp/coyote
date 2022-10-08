package cmd

import (
	"github.com/lnsp/coyote/tracker"

	"github.com/spf13/cobra"
)

var (
	trackOutput  string
	trackTaverns []string
)

var trackCmd = &cobra.Command{
	Use:   "track [file]",
	Short: "Create a file tracker with a specified host tavern",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		trackPath := args[0]
		if len(trackOutput) == 0 {
			trackOutput = trackPath + ".coyote"
		}
		// Use progress bar to track updates
		t, err := tracker.Track(trackTaverns, trackPath)
		if err != nil {
			return err
		}
		return tracker.Save(t, trackOutput)
	},
}

func init() {
	trackCmd.Flags().StringArrayVarP(&trackTaverns, "taverns", "s", []string{"https://localhost:6443"}, "Tavern hosts which serves the peer sessions")
	trackCmd.Flags().StringVarP(&trackOutput, "output", "o", "", "Where to put the Tracker")
	rootCmd.AddCommand(trackCmd)
}
