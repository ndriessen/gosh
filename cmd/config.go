package cmd

import "github.com/spf13/cobra"

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Displays information about configuration options for Gosh",
		Long: `Configuration options:
1) ~/.gosh/config.yml or ./.gosh/config.yml

You can configure different ways to authenticate for GIT operations
NOTE: only use one of these in your actual config

1.1) Basic Auth
Auth:
  Type: Basic
  User: your-user
  Pass: your-pass-base64-encoded

IMPORTANT: Use HTTPS URLs for your GIT repository URLs when using basic auth!

1.2) SSH Key
Auth:
  Type: ssh
  Private_Key_File: ~/.ssh/id_rsa
  Private_Key_Pass: your-private-key-pass-base64-encoded

If you want to use Gosh artifacts and replacements, you also have to configure these

ArtifactRepositories:
  maven:
    default: https://your.maven.repo/repository/tm-released
    STAGE_NAME: https://your.maven.repo/repository/tm-tested
    OTHER_STAGE: https://your.maven.repo/repository/tm-published
  docker:
    default: "your.docker.registry"
  ...

You can add as many as you want/need

2) Using ENV variables
You can also configure authentication using ENV variables

2.1) Basic Auth
GOSH_AUTH_TYPE=basic
GOSH_AUTH_USER=username
GOSH_AUTH_PASS=your-pass-base64-encoded

2.2) SSH
GOSH_AUTH_TYPE=ssh
GOSH_AUTH_PRIVATE_KEY_FILE=~/.ssh/id_rsa
GOSH_AUTH_PRIVATE_KEY_PASS=your-private-key-pass-base64-encoded

`,
	}
)

func init() {
	rootCmd.AddCommand(configCmd)
}
