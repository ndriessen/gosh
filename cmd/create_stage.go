package cmd

import (
	"github.com/spf13/cobra"
	"gosh/gitops"
	"log"
)

var (
	createStageCmd = &cobra.Command{
		Use:  "stage [NAME]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stageName := GetArg(args, 0)
			stage := gitops.NewStage(stageName)
			if err := stage.Create(); err != nil {
				log.Fatalln("Error creating stage", stageName, err)
			}
		},
	}
)

func init() {
	createCmd.AddCommand(createStageCmd)
}
