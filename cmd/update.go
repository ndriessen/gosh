package cmd

import "github.com/spf13/cobra"

var (
	updateCmd = &cobra.Command{
		Use:  "update",
		Args: cobra.ExactArgs(1),
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
}
