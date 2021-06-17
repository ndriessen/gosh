package cmd

import (
	"github.com/spf13/cobra"
	"gosh/gitops"
	"gosh/log"
)

const (
	fromStageFlag   = "from-stage"
	fromReleaseFlag = "from-release"
)

var (
	createReleaseCmd = &cobra.Command{
		Use:  "release [prefix/NAME]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			releaseName := GetArg(args, 0)
			if release, err := gitops.NewReleaseFromFullName(releaseName); err == nil {
				if flag, value, err := GetMutuallyExclusiveStringFlag(cmd, fromStageFlag, fromReleaseFlag); err == nil {
					switch flag {
					case fromStageFlag:
						if err = release.CreateFromStage(value); err != nil {
							log.Fatal(err, "Error creating release %s from stage %s", releaseName, value)
						}
					case fromReleaseFlag:
						if err = release.CreateFromRelease(value); err != nil {
							log.Fatal(err, "Error creating release %s from release %s", releaseName, value)
						}
					}
				} else {
					log.Fatal(err, "Please specify either --from-stage or --from-release")
				}
			} else {
				log.Fatal(err, "Error creating release %s", releaseName)
			}
		},
	}
)

func init() {
	createReleaseCmd.Flags().StringP(fromStageFlag, "S", "", "--from-stage|-S STAGE")
	createReleaseCmd.Flags().StringP(fromReleaseFlag, "R", "", "--from-release|-R PREFIX/RELEASE_NAME")
	createCmd.AddCommand(createReleaseCmd)
}
