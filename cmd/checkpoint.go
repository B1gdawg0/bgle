package cmd

import (
	"bgle/config"
	"bgle/utils"

	"github.com/spf13/cobra"
)

var checkpointCmd = &cobra.Command{
    Use:   "checkpoint [name:profile]",
    Short: "Create a checkpoint (profile) for a project",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name, profile, err := utils.ParseNameProfile(args[0])
        if err != nil {
            utils.PrintError("Error: "+err.Error())
        }

        cfg := config.InteractiveProfileCreation(name, profile)

        err = config.SaveProfile(cfg)
        if err != nil {
            utils.PrintError("Failed to save profile: "+ err.Error())
        }
    },
}

func init() {
    rootCmd.AddCommand(checkpointCmd)
}

