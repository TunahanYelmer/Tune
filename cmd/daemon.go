package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/config"
	"github.com/tunahanyelmer/Tune/internal/daemon"
	"github.com/tunahanyelmer/Tune/internal/providers/registry"
)

var daemonCmd = &cobra.Command{
	Use:    "__daemon",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		provider, err := registry.Load(cfg.Provider)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "tune daemon started, provider:", cfg.Provider)
		return daemon.Serve(provider)
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}