package cmd

import (
	"github.com/spf13/cobra"
	"gosh/log"
)

var (
	listVersionsCmd = &cobra.Command{
		Use: "versions [FLAGS] [APP_NAME]",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("list versions")
		},
	}
)

func init() {
	listCmd.AddCommand(listVersionsCmd)
}
