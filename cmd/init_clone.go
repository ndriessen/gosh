package cmd

import (
	"github.com/spf13/cobra"
	"gosh/git"
	"gosh/log"
)

var (
	initCloneCmd = &cobra.Command{
		Use:  "clone",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := GetArg(args, 0)
			if _, err := git.NewDeploymentRepository(url, true); err != nil {
				log.Fatal(err, "Error cloning deployment repository in working directory")
			}
		},
	}
)

func init() {
	initCmd.AddCommand(initCloneCmd)
}
