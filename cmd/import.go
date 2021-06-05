package cmd

import (
	"github.com/spf13/cobra"
	gosh_import "gosh/import"
	"gosh/log"
)

var (
	importCmd = &cobra.Command{
		Use:  "import PLUGIN",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := GetArg(args, 0)
			if err := gosh_import.Import(name); err != nil {
				log.Fatal(err, "error running import with plugin %s", name)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(importCmd)
}
