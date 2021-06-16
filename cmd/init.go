package cmd

import "github.com/spf13/cobra"

var (
	initCmd = &cobra.Command{
		Use: "init",
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}
