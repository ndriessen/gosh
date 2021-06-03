package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gosh/list"
	"gosh/log"
)

var (
	listVersionsCmd = &cobra.Command{
		Use:  "versions {--stage STAGE | --release RELEASE} [FLAGS]... [APP_NAME]",
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
				if data, err := list.Render(
					GetStringFlag(cmd, OutputFlag, list.DefaultOutputFormat),
					l.GetVersions(GetStringFlag(cmd, GroupFlag, ""), GetArg(args, 0)),
				); err == nil {
					fmt.Println(data)
				} else {
					log.Fatal(err, "Could not list versions")
				}
			} else {
				log.Fatal(err, "Could not list versions")
			}
		},
	}
)

func init() {
	AddStageFlag(listVersionsCmd)
	AddReleaseFlag(listVersionsCmd)
	AddGroupFlag(listVersionsCmd)
	AddOutputFlag(listVersionsCmd)
	listCmd.AddCommand(listVersionsCmd)
}
