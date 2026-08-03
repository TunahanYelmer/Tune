package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/daemon"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the currently playing track",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.EnsureRunning(); err != nil {
			return err
		}

		resp, err := daemon.Send(daemon.Request{Action: daemon.ActionCurrent})
		if err != nil {
			return err
		}
		if !resp.OK {
			return fmt.Errorf("%s", resp.Error)
		}

		fmt.Println("Now Playing")
		fmt.Printf("🎵 %s\n👤 %s\n💿 %s\n", resp.Track.Title, resp.Track.Artist, resp.Track.Album)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}