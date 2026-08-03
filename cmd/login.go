package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/daemon"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the configured provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.EnsureRunning(); err != nil {
			return err
		}

		resp, err := daemon.Send(daemon.Request{Action: daemon.ActionLogin})
		if err != nil {
			return err
		}
		if !resp.OK {
			return fmt.Errorf("%s", resp.Error)
		}

		fmt.Println("✅ Logged in")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}