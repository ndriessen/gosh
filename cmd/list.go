package cmd

import (
	"github.com/spf13/cobra"
	"gosh/gitops"
	"gosh/log"
)

var (
	listCmd = &cobra.Command{
		Use: "list",
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func LoadVersionsList(flag string, value string) (gitops.AppList, error) {
	log.Debugf("Loading %s versions for %s", flag, value)
	switch flag {
	case StageFlag:
		{
			stage := gitops.NewStage(value)
			if err := stage.Read(); err == nil {
				return stage, nil
			} else {
				return nil, log.Errf(err, "Error loading stage %s", value)
			}

		}
	case ReleaseFlag:
		{
			if release, err := gitops.NewReleaseFromFullName(value); err == nil {
				if err = release.Read(); err == nil {
					return release, nil
				} else {
					return nil, log.Errf(err, "Error loading release %s", value)
				}
			} else {
				return nil, log.Errf(err, "Error loading release %s", value)
			}
		}
	default:
		return nil, log.Errf(gitops.ResourceDoesNotExistErr, "unknown version list type %s", flag)
	}

}
