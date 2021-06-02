package cmd

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"gosh/log"
	"gosh/util"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gosh",
		Short: "Gosh or GitOps Shell offer convenience for interacting with a deployment repository based on GitOps concepts.",
		PreRun: func(cmd *cobra.Command, args []string) {

		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose mode to output logging (default: false)")
	rootCmd.PersistentFlags().BoolP("trace", "t", false, "enable trace logging, only needed for development/testing (default: false)")
	rootCmd.PersistentFlags().StringP("workdir", "w", "$PWD", "specify the working directory for gosh (default: $PWD)")

	cobra.OnInitialize(handleGlobalFlags)
}

func GetStringFlag(cmd *cobra.Command, name string, defaultValue string) string {
	if value, err := cmd.Flags().GetString(name); err == nil {
		return value
	} else if value, err = cmd.PersistentFlags().GetString(name); err == nil {
		return value
	}
	return defaultValue
}

func GetBoolFlag(cmd *cobra.Command, name string, defaultValue bool) bool {
	if value, err := cmd.Flags().GetBool(name); err == nil {
		return value
	} else if value, err = cmd.PersistentFlags().GetBool(name); err == nil {
		return value
	}
	return defaultValue
}

func Execute() error {
	err := rootCmd.Execute()
	return log.CheckErr(err)
}

func handleGlobalFlags() {
	if wd := GetStringFlag(rootCmd, "workdir", ""); wd != "" {
		log.Debug("Setting workdir from flag", wd)
		util.Context.WorkingDir = os.ExpandEnv(wd)
	}
	if verbose := GetBoolFlag(rootCmd, "verbose", false); verbose {
		log.Debug("Enabling verbose logging")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if trace := GetBoolFlag(rootCmd, "trace", false); trace {
		log.Debug("Enabling tracing")
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}
