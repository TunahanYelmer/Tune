package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/daemon"
)

var playCmd = &cobra.Command{
	Use:   "play <song>",
	Short: "Play a song",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.EnsureRunning(); err != nil {
			return err
		}
		query := strings.TrimSpace(strings.Join(args, " "))
		resp, err := daemon.Send(daemon.Request{Action: daemon.ActionPlay, Query: query})
		if err != nil {
			return err
		}
		if !resp.OK {
			return fmt.Errorf("%s", resp.Error)
		}
		fmt.Printf("▶ Playing: %s\n", query)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
}