package cmd

import (
	"github.com/spf13/cobra"
	"gosh/git"
)

var (
	initCloneCmd = &cobra.Command{
		Use: "clone",
		Run: func(cmd *cobra.Command, args []string) {
			git.InitializeGit(true)
		},
	}
)

func init() {
	initCmd.AddCommand(initCloneCmd)
}
