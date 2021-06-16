package cmd

import (
	"github.com/spf13/cobra"
	"gosh/gitops"
	"log"
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
							log.Fatalf("Error creating release %s from stage %s (error: %v)", releaseName, value, err)
						}
					case fromReleaseFlag:
						if err = release.CreateFromRelease(value); err != nil {
							log.Fatalf("Error creating release %s from release %s (error: %v)", releaseName, value, err)
						}
					}
				}
			} else {
				log.Fatalln("Error creating release", releaseName, err)
			}
		},
	}
)

func init() {
	createReleaseCmd.Flags().StringP(fromStageFlag, "S", "", "--from-stage|-S STAGE")
	createReleaseCmd.Flags().StringP(fromReleaseFlag, "R", "", "--from-release|-R PREFIX/RELEASE_NAME")
	createCmd.AddCommand(createReleaseCmd)
}
