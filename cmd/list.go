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

func LoadAppList(appListType string, appListName string) (gitops.AppList, error) {
	log.Debugf("Loading %s versions for %s", appListType, appListName)
	switch appListType {
	case StageFlag:
		{
			stage := gitops.NewStage(appListName)
			if err := stage.Read(); err == nil {
				return stage, nil
			} else {
				return nil, err
			}

		}
	case ReleaseFlag:
		{
			if release, err := gitops.NewReleaseFromFullName(appListName); err == nil {
				if err = release.Read(); err == nil {
					return release, nil
				} else {
					return nil, log.Errf(err, "Error loading release %s", appListName)
				}
			} else {
				return nil, err
			}
		}
	default:
		return nil, log.Errf(gitops.ResourceDoesNotExistErr, "unknown version list type %s", appListType)
	}

}
