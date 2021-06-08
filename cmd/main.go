package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	StageFlag   = "stage"
	ReleaseFlag = "release"
	GroupFlag   = "group"
	OutputFlag  = "output"
)

var RequiredFlagMissingErr = errors.New("required flag is missing")

func GetArg(args []string, position int) string {
	if len(args) > position {
		return args[position]
	}
	return ""
}

func GetRequiredArg(args []string, position int) (string, error) {
	if len(args) > position {
		return args[position], nil
	}
	return "", RequiredFlagMissingErr
}

func AddStageFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(StageFlag, "s", "", "--stage|-s STAGE")
}

func AddReleaseFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(ReleaseFlag, "r", "", "--release|-r RELEASE")
}

func AddGroupFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(GroupFlag, "g", "", "--group|-g GROUP")
}

func AddOutputFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(OutputFlag, "o", "yaml", "--output|-o yaml|properties (default: yaml)")
}

func GetStringFlag(cmd *cobra.Command, name string, defaultValue string) string {
	if value, err := cmd.Flags().GetString(name); err == nil && value != "" {
		return value
	} else if value, err = cmd.PersistentFlags().GetString(name); err == nil && value != "" {
		return value
	}
	return defaultValue
}

var (
	RequiredFlagNotSetErr        = errors.New("required flag not set")
	MutuallyExclusiveFlagsSetErr = errors.New("mutually exclusive flags set")
)

func GetMutuallyExclusiveStringFlag(cmd *cobra.Command, flags ...string) (flag string, value string, err error) {
	valuesSet := 0
	for _, name := range flags {
		if val, err := cmd.Flags().GetString(name); err == nil && strings.TrimSpace(val) != "" {
			valuesSet++
			flag = name
			value = val
		}
	}
	switch valuesSet {
	case 0:
		return "", "", RequiredFlagNotSetErr
	case 1:
		return flag, value, nil
	default:
		return "", "", MutuallyExclusiveFlagsSetErr
	}
}

func GetBoolFlag(cmd *cobra.Command, name string, defaultValue bool) bool {
	if value, err := cmd.Flags().GetBool(name); err == nil {
		return value
	} else if value, err = cmd.PersistentFlags().GetBool(name); err == nil {
		return value
	}
	return defaultValue
}
