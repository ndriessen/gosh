package cmd

import (
	"github.com/spf13/cobra"
	"gosh/gitops"
	"log"
)

var (
	createAppCmd = &cobra.Command{
		Use:  "app",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appName := GetArg(args, 0)
			appGroupName := GetStringFlag(cmd, GroupFlag, "")
			templateName := GetStringFlag(cmd, TemplateFlag, "")
			appGroup := gitops.NewAppGroup(appGroupName)
			app := gitops.NewApp(appName, appGroup)
			if err := app.CreateFromTemplate(templateName); err != nil {
				log.Fatalln("Error creating app", err)
			}
		},
	}
)

func init() {
	AddGroupFlag(createAppCmd)
	AddTemplateFlag(createAppCmd)
	_ = createAppCmd.MarkFlagRequired(GroupFlag)
	createCmd.AddCommand(createAppCmd)
}
