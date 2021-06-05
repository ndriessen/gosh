package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	gosh_import "gosh/import"
	"gosh/log"
)

var (
	cmdImportPlugins = &cobra.Command{
		Use:   "plugins",
		Short: "Lists all available import plugins",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Bundled import plugins:")
			for k := range gosh_import.BundledPlugins {
				fmt.Println("-", k)
			}
			fmt.Println()
			fmt.Println("Available import plugins:")
			if plugins, err := gosh_import.ListPlugins(); err == nil {
				if len(plugins) == 0 {
					fmt.Println("- NO PLUGINS INSTALLED")
				}
				for _, p := range plugins {
					fmt.Println("-", p)
				}
				fmt.Println("\nTo install a plugin, copy the file to ~/.gosh/plugins")
			} else {
				log.Fatal(err, "Could not list available plugins")
			}
		},
	}
)

func init() {
	importCmd.AddCommand(cmdImportPlugins)
}
