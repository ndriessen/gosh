package util

import (
	"github.com/spf13/viper"
	"gosh/log"
	"os"
	"path/filepath"
)

const (
	GoshConfigDir         = ".gosh"
	GoshDefaultConfigFile = "$HOME/" + GoshConfigDir + "/" + GoshConfigFile
	GoshConfigFile        = "config.yml"
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
	viper.SetEnvPrefix("GOSH")
	viper.AutomaticEnv()
	viper.SetConfigFile(os.ExpandEnv(GoshDefaultConfigFile))
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
