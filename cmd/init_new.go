package cmd

import (
	"github.com/spf13/cobra"
	"gosh/log"
)

var (
	initNewCmd = &cobra.Command{
		Use:   "new GIT_REPOSITORY_URL",
		Short: "Initializes an empty deployment repository in the specified GIT repository (the repository has to exists and you need to configure SSH keys",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Debugf("init new")
		},
	}
)

func init() {
	initCmd.AddCommand(initNewCmd)
}
