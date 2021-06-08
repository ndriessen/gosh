package cmd

import (
	"github.com/spf13/cobra"
	"gosh/log"
)

var (
	updateVersionCmd = &cobra.Command{
		Use:  "version {--stage STAGE|--release RELEASE} APP VERSION",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			appName := GetArg(args, 0)
			version := GetArg(args, 1)
			flag, value, err := GetMutuallyExclusiveStringFlag(cmd, "stage", "release")
			if err == MutuallyExclusiveFlagsSetErr {
				log.Fatal(err, "You must specify either --stage or --release, not both")
			}
			if err == RequiredFlagNotSetErr {
				log.Fatal(err, "You must specify --stage or --release")
			}
			if appList, err := LoadAppList(flag, value); err == nil {
				if err = appList.UpdateVersion(appName, version); err == nil {
					log.Infof("Updated app %s to version %s for %s %s", appName, version, flag, value)
				} else {
					log.Fatal(err, "Error updating app %s to version %s for %s %s", appName, version, flag, value)
				}
			} else {
				log.Fatal(err, "Error loading %s %s", flag, value)
			}
		},
	}
)

func init() {
	AddReleaseFlag(updateVersionCmd)
	AddStageFlag(updateVersionCmd)
	updateCmd.AddCommand(updateVersionCmd)
}
