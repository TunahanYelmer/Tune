package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/daemon"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for a track",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.EnsureRunning(); err != nil {
			return err
		}

		query := strings.Join(args, " ")
		resp, err := daemon.Send(daemon.Request{Action: daemon.ActionSearch, Query: query})
		if err != nil {
			return err
		}
		if !resp.OK {
			return fmt.Errorf("%s", resp.Error)
		}

		for i, t := range resp.Tracks {
			fmt.Printf("%d. %s - %s\n", i+1, t.Artist, t.Title)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}