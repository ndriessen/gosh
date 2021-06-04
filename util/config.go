package util

import (
	"github.com/spf13/viper"
	"gosh/log"
	"os"
	"path/filepath"
)

const (
	GOSH_CONFIG_DIR          = ".gosh"
	GOSH_DEFAULT_CONFIG_FILE = "$HOME/" + GOSH_CONFIG_DIR + "/" + GOSH_CONFIG_FILE
	GOSH_CONFIG_FILE         = "config.yml"
)

type DeploymentRepository struct {
	Url               string
	SshKey            string
	SshPrivateKeyPass string
}

func newDeploymentRepository(settings map[string]interface{}) DeploymentRepository {
	return DeploymentRepository{
		Url:               settings["url"].(string),
		SshKey:            settings["sshkey"].(string),
		SshPrivateKeyPass: settings["sshprivatekeypass"].(string),
	}
}

type GoshConfig struct {
	DeploymentRepository
	ArtifactRepositories map[string]map[string]string
}

var Config = &GoshConfig{}

func InitializeConfig() {
	v := viper.New()
	projectConfigFile := filepath.Join(Context.WorkingDir, GOSH_CONFIG_DIR, GOSH_CONFIG_FILE)
	v.SetConfigFile(projectConfigFile)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err, "Could not read configuration")
		}
	}
	viper.SetEnvPrefix("GOSH")
	viper.AutomaticEnv()
	viper.SetConfigFile(os.ExpandEnv(GOSH_DEFAULT_CONFIG_FILE))
	if err := viper.ReadInConfig(); err == nil {
		if err = viper.MergeConfigMap(v.AllSettings()); err != nil {
			log.Fatal(err, "Error merging configuration")
		}
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err, "Could not read configuration")
		}
	}
	Config.DeploymentRepository = newDeploymentRepository(viper.Get("deploymentrepository").(map[string]interface{}))
	Config.ArtifactRepositories = make(map[string]map[string]string, 0)
	settings := viper.Get("artifactrepositories").(map[string]interface{})
	for name, urls := range settings {
		result := map[string]string{}
		for k, v := range urls.(map[string]interface{}) {
			result[k] = v.(string)
		}
		Config.ArtifactRepositories[name] = result
	}
	log.Debugf("Loaded configuration %+v", Config)
}
