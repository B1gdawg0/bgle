package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "bgle",
    Short: "bgle is a CLI pipeline profile tool",
}

func Execute() {
    cobra.CheckErr(rootCmd.Execute())
}

func init() {
    rootCmd.AddCommand(checkpointCmd)
    rootCmd.AddCommand(useCmd)
}