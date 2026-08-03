package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/config"
	"github.com/tunahanyelmer/Tune/internal/daemon"
	"github.com/tunahanyelmer/Tune/internal/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive terminal UI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.EnsureRunning(); err != nil {
			return fmt.Errorf("starting daemon: %w", err)
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		p := tea.NewProgram(tui.New(cfg.Provider))
		_, err = p.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}