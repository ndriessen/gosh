package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gosh/list"
	"gosh/log"
)

var (
	listArtifactsCmd = &cobra.Command{
		Use:  "artifacts {--stage STAGE | --release RELEASE} [FLAGS]... [APP_NAME]",
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Tracef("running command list versions with args: %v", args)
			flag, value, err := GetMutuallyExclusiveStringFlag(cmd, "stage", "release")
			if err == MutuallyExclusiveFlagsSetErr {
				log.Fatal(err, "You must specify either --stage or --release, not both")
			}
			if err == RequiredFlagNotSetErr {
				log.Fatal(err, "You must specify --stage or --release")
			}
			if l, err := LoadVersionsList(flag, value); err == nil {
				if artifacts, err := l.GetArtifacts(GetStringFlag(cmd, GroupFlag, ""), GetArg(args, 0), "maven"); err == nil {
					if data, err := list.Render(GetStringFlag(cmd, OutputFlag, list.DefaultOutputFormat), artifacts); err == nil {
						fmt.Println(data)
					} else {
						log.Fatal(err, "Could not list versions")
					}
				} else {
					log.Fatal(err, "Could not list artifacts, make sure all apps have artifacts defined")
				}
			} else {
				log.Fatal(err, "Could not list versions")
			}
		},
	}
)

func init() {
	AddStageFlag(listArtifactsCmd)
	AddReleaseFlag(listArtifactsCmd)
	AddGroupFlag(listArtifactsCmd)
	AddOutputFlag(listArtifactsCmd)
	listCmd.AddCommand(listArtifactsCmd)
}
