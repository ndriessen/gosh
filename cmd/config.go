package cmd

import "github.com/spf13/cobra"

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Displays information about configuration options for Gosh",
		Long: `Configuration options:

Gosh reads configuration from 3 places:
- '~/.gosh/config.yml' in your home dir
- './.gosh/config.yml' in your working dir (project specific)
- from ENV variables

1) Authentication for GIT repositories

1.1) Basic Auth
1.1.1) In config files
Auth:
  Type: Basic
  User: your-user
  Pass: your-pass-base64-encoded
1.1.2) Using ENV
GOSH_AUTH_TYPE=basic
GOSH_AUTH_USER=username
GOSH_AUTH_PASS=your-pass-base64-encoded

IMPORTANT: Use HTTPS URLs for your GIT repository URLs when using basic auth!

1.2) SSH Key
1.2.1) In config files
Auth:
  Type: ssh
  Private_Key_File: ~/.ssh/id_rsa
  Private_Key_Pass: your-private-key-pass-base64-encoded
1.1.2) Using ENV
GOSH_AUTH_TYPE=ssh
GOSH_AUTH_PRIVATE_KEY_FILE=~/.ssh/id_rsa
GOSH_AUTH_PRIVATE_KEY_PASS=your-private-key-pass-base64-encoded

2) Output configuration
You can set some configuration that is used by commands that output lists (like list versions and list artifacts)
Suffixes are optional, and default to an empty string, output format can also be specified as a command flag, if set
this default will be used when no flag is passed

2.1) In config files
Output:
  Default_Format: yaml|properties
  Versions_Key_Suffix: "version" # will yield  [APP_NAME].version=[APP_VERSION] for list versions
  Artifacts_Key_Suffix: "" # will yield  [APP_NAME]=[APP_ARTIFACT] for list versions
2.2) Using ENV
GOSH_OUTPUT_DEFAULT_FORMAT=yaml|properties
GOSH_OUTPUT_VERSIONS_KEY_SUFFIX=version
GOSH_OUTPUT_ARTIFACTS_KEY_SUFFIX=

3) Artifact repositories
If you want to use Gosh artifacts and replacements, you also have to configure these
3.1) In config files
ArtifactRepositories:
  maven:
    default: https://your.maven.repo/repository/tm-released
    STAGE_NAME: https://your.maven.repo/repository/tm-tested
    OTHER_STAGE: https://your.maven.repo/repository/tm-published
  docker:
    default: "your.docker.registry"
  ...

You can add as many as you want/need

`,
	}
)

func init() {
	rootCmd.AddCommand(configCmd)
}
