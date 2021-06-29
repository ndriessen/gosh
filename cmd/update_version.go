package cmd

import (
	"github.com/spf13/cobra"
	"gosh/git"
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
				if repo, err := git.NewDeploymentRepository("", false); err == nil {
					if err = appList.UpdateVersion(appName, version); err == nil {
						push := GetBoolFlag(cmd, "push", false)
						if push {
							err = repo.Push(GetStringFlag(cmd, "message", ""))
						}
						if err == nil {
							log.Infof("Updated app %s to version %s for %s %s", appName, version, flag, value)
						} else {
							log.Fatal(err, "Error pushing updates to deployment repository")
						}
					} else {
						log.Fatal(err, "Error updating app %s to version %s for %s %s", appName, version, flag, value)
					}
				} else {
					log.Fatal(err, "Error opening working dir as Git repository")
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
	updateVersionCmd.Flags().BoolP("push", "p", false, "--push|-p   Push changes to the remote repository (default: false)")
	updateVersionCmd.Flags().StringP("message", "m", "", "--message|-m \"COMMIT MESSAGE\" (optional, only used when --push is specified)")
	updateCmd.AddCommand(updateVersionCmd)
}
