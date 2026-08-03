package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/daemon"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause or resume playback",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.EnsureRunning(); err != nil {
			return err
		}

		resp, err := daemon.Send(daemon.Request{Action: daemon.ActionPause})
		if err != nil {
			return err
		}
		if !resp.OK {
			return fmt.Errorf("%s", resp.Error)
		}

		fmt.Println("⏯ Toggled play/pause")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}