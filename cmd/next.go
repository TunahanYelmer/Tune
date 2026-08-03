package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/daemon"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Skip to the next track",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.EnsureRunning(); err != nil {
			return err
		}

		resp, err := daemon.Send(daemon.Request{Action: daemon.ActionNext})
		if err != nil {
			return err
		}
		if !resp.OK {
			return fmt.Errorf("%s", resp.Error)
		}

		fmt.Println("⏭ Skipped to next")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}