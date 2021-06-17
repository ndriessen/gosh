package cmd

import (
	"github.com/spf13/cobra"
	"gosh/git"
)

var (
	initCloneCmd = &cobra.Command{
		Use:  "clone",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			//url := GetArg(args,0)
			git.InitializeGit(true)
		},
	}
)

func init() {
	initCmd.AddCommand(initCloneCmd)
}
