package cmd

import (
	"github.com/spf13/cobra"
	"gosh/log"
)

var (
	initNewCmd = &cobra.Command{
		Use: "new",
		Run: func(cmd *cobra.Command, args []string) {
			log.Debugf("init new")
		},
	}
)

func init() {
	initCmd.AddCommand(initNewCmd)
}
