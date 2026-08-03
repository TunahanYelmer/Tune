package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/config"
)

var providerCmd = &cobra.Command{
	Use:   "provider [name]",
	Short: "Show or set the active music provider",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			fmt.Println("Current provider:", cfg.Provider)
			return nil
		}

		cfg.Provider = args[0]
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("✅ Provider set to %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(providerCmd)
}