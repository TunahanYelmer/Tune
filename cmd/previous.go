package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/daemon"
)

var previousCmd = &cobra.Command{
	Use:   "previous",
	Short: "Go back to the previous track",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.EnsureRunning(); err != nil {
			return err
		}

		resp, err := daemon.Send(daemon.Request{Action: daemon.ActionPrev})
		if err != nil {
			return err
		}
		if !resp.OK {
			return fmt.Errorf("%s", resp.Error)
		}

		fmt.Println("⏮ Went to previous")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(previousCmd)
}