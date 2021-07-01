package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gosh/util"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gosh " + util.Context.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
