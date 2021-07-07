package util

import (
	"errors"
	"github.com/spf13/viper"
	"gosh/log"
	"os"
	"path/filepath"
	"strings"
)

const (
	GoshConfigDir  = ".gosh"
	GoshConfigFile = "config.yml"
)

type AuthType int

const (
	BasicAuth AuthType = iota + 1
	SshKey
)

var UnsupportedGitAuthTypeErr = errors.New("unsupported git auth type")

func newAuthType(value string) (AuthType, error) {
	switch strings.ToLower(value) {
	case "basic":
		return BasicAuth, nil
	case "ssh":
		return SshKey, nil
	}
	return 0, UnsupportedGitAuthTypeErr
}

type AuthConfig interface {
	Type() AuthType
}

type BasicAuthConfig struct {
	User string
	Pass string
}

func (auth BasicAuthConfig) Type() AuthType {
	return BasicAuth
}

type SshAuthConfig struct {
	PrivateKeyFile string
	PrivateKeyPass string
}

func (auth SshAuthConfig) Type() AuthType {
	return SshKey
}

func newGitAuthConfig(authConfig *authConfigDef) (AuthConfig, error) {
	t, err := newAuthType(authConfig.Type)
	if err != nil {
		return nil, err
	}
	switch t {
	case BasicAuth:
		return BasicAuthConfig{
			User: authConfig.User,
			Pass: authConfig.Pass,
		}, nil
	case SshKey:
		return SshAuthConfig{
			PrivateKeyFile: authConfig.PrivateKeyFile,
			PrivateKeyPass: authConfig.PrivateKeyPass,
		}, nil
	default:
		return nil, UnsupportedGitAuthTypeErr
	}
}

type GoshConfig struct {
	Auth                 AuthConfig
	ArtifactRepositories map[string]map[string]string
}

type authConfigDef struct {
	Type           string
	User           string
	Pass           string
	PrivateKeyFile string `mapstructure:"private_key_file"`
	PrivateKeyPass string `mapstructure:"private_key_pass"`
}

var Config = &GoshConfig{}

func InitializeConfig() {
	v := viper.New()
	projectConfigFile := filepath.Join(Context.WorkingDir, GoshConfigDir, GoshConfigFile)
	if _, err := os.Stat(projectConfigFile); err == nil {
		v.SetConfigFile(projectConfigFile)
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Fatal(err, "Could not read configuration")
			}
		}
	} else {
		log.Debugf("no project specific config file found in %s, skipping...", Context.WorkingDir)
	}
	vpr := viper.New()
	vpr.SetEnvPrefix("GOSH")
	vpr.AutomaticEnv()
	vpr.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	homedir, _ := os.UserHomeDir()
	configFile := filepath.Join(homedir, GoshConfigDir, GoshConfigFile)
	vpr.SetConfigFile(configFile)
	if _, err := os.Stat(configFile); err == nil {
		if err := vpr.ReadInConfig(); err == nil {
			if err = vpr.MergeConfigMap(v.AllSettings()); err != nil {
				log.Fatal(err, "Error merging configuration")
			}
		} else {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Fatal(err, "Could not read configuration")
			}
		}
	}
	//for backward compat
	if vpr.IsSet("deploymentrepository") {
		authConfig := &SshAuthConfig{
			PrivateKeyFile: os.ExpandEnv(vpr.GetString("deploymentrepository.sshkey")),
			PrivateKeyPass: vpr.GetString("deploymentrepository.sshprivatekeypass"),
		}
		Config.Auth = authConfig
	}
	//new config properties override the deprecated ones
	if vpr.IsSet("auth.type") {
		authConfig := &authConfigDef{
			Type:           vpr.GetString("auth.type"),
			User:           vpr.GetString("auth.user"),
			Pass:           vpr.GetString("auth.pass"),
			PrivateKeyFile: vpr.GetString("auth.private_key_file"),
			PrivateKeyPass: vpr.GetString("auth.private_key_pass"),
		}
		auth, err := newGitAuthConfig(authConfig)
		if err != nil {
			log.Fatal(err, "Invalid auth configuration %+v", authConfig)
		}
		Config.Auth = auth
		log.Debugf("Using %s auth config", authConfig.Type)
	}
	Config.ArtifactRepositories = make(map[string]map[string]string, 0)
	if vpr.IsSet("artifactrepositories") {
		settings := vpr.Get("artifactrepositories").(map[string]interface{})
		for name, urls := range settings {
			result := map[string]string{}
			for k, v := range urls.(map[string]interface{}) {
				result[k] = v.(string)
			}
			Config.ArtifactRepositories[name] = result
		}
	}
	log.Debugf("Loaded configuration %+v", Config)
}
