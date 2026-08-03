package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/config"
	"github.com/tunahanyelmer/Tune/internal/providers/registry"
	"github.com/tunahanyelmer/Tune/internal/service"
)

var Music *service.MusicService

var rootCmd = &cobra.Command{
	Use:   "tune",
	Short: "Control your music from the terminal",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initMusic)
}

func initMusic() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading config: %v\n", err)
		os.Exit(1)
	}

	provider, err := registry.Load(cfg.Provider)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading provider: %v\n", err)
		os.Exit(1)
	}

	Music = service.NewMusicService(provider)
}