package cmd

import (
	"github.com/spf13/cobra"
	gosh_import "gosh/import"
	"gosh/log"
	"strings"
)

const (
	all = "all"
)

var (
	importCmd = &cobra.Command{
		Use:  "import PLUGIN [all|apps|stages|releases]... (default=all)",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := GetArg(args, 0)
			var argList string
			if len(args) > 1 {
				argList = strings.ToLower(strings.Join(args[1:], "|"))
			} else {
				argList = all
			}
			apps := strings.Contains(argList, "apps") || strings.Contains(argList, all)
			stages := strings.Contains(argList, "stages") || strings.Contains(argList, all)
			releases := strings.Contains(argList, "releases") || strings.Contains(argList, all)
			if err := gosh_import.Import(name, apps, stages, releases); err != nil {
				log.Fatal(err, "error running import with plugin %s", name)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(importCmd)
}
