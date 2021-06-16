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
		//PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//	log.Debug("Pre-run init")
		//	if cmd != initCmd && cmd != initNewCmd {
		//		git.InitializeGit(false)
		//	}
		//},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose mode to output logging (default: false)")
	rootCmd.PersistentFlags().BoolP("trace", "t", false, "enable trace logging, only needed for development/testing (default: false)")
	rootCmd.PersistentFlags().StringP("workdir", "w", "", "specify the working directory for gosh (default: $PWD)")

	cobra.OnInitialize(handleGlobalFlags)
	cobra.OnInitialize(util.InitializeConfig)
}

func Execute() error {
	err := rootCmd.Execute()
	return log.CheckErr(err)
}

func handleGlobalFlags() {
	zerolog.SetGlobalLevel(zerolog.FatalLevel)
	if verbose := GetBoolFlag(rootCmd, "verbose", false); verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.SetDetailedLogging()
		log.Debug("Enabling verbose logging")
	}
	if trace := GetBoolFlag(rootCmd, "trace", false); trace {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.SetDetailedLogging()
		log.Tracef("Enabling tracing")
	}
	if wd := GetStringFlag(rootCmd, "workdir", ""); wd != "" {
		log.Debugf("Setting workdir from flag: %s", os.ExpandEnv(wd))
		util.Context.WorkingDir = os.ExpandEnv(wd)
	}
}
