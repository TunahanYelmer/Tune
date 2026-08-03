package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/tunahanyelmer/Tune/internal/daemon"
)

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "change the volume",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.EnsureRunning(); err != nil {
			return err
		}
		
		

			level, err := strconv.Atoi(args[0])
			if err != nil {
			fmt.Println("Level must be a number")
			return err
}


		

		
		resp, err := daemon.Send(daemon.Request{Action: daemon.ActionVolume, Level:level })
		if err != nil {
			return err
		}
		if !resp.OK {
			return fmt.Errorf("%s", resp.Error)
		}
		fmt.Printf("▶ Playing: %d\n", level)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(volumeCmd)
}