package cmd

import (
	"github.com/spf13/cobra"
	"gosh/git"
	"gosh/log"
)

var (
	initNewCmd = &cobra.Command{
		Use:   "new [GIT_REPOSITORY_URL]",
		Short: "Initializes an empty deployment repository, if you specify a repository URL, it will first be cloned",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := GetArg(args, 0)
			if url != "" {
				if repo, err := git.NewDeploymentRepository(url, true); err != nil {
					log.Fatal(err, "Unable to initialize deployment repository in working directory")
				} else {
					if err = repo.InitFromTemplate(); err != nil {
						log.Fatal(err, "Unable to initialize deployment repository from template in working directory")
					}
				}
			}
		},
	}
)

func init() {
	initCmd.AddCommand(initNewCmd)
}
