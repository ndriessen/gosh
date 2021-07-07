package util

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
	homedir string
	workdir string
}

func (suite *ConfigTestSuite) SetupSuite() {
}

func (suite *ConfigTestSuite) SetupTest() {
	os.Clearenv()
	suite.homedir = filet.TmpDir(suite.T(), "")
	suite.workdir = filet.TmpDir(suite.T(), "")
	_ = os.Setenv("HOME", suite.homedir)
	_ = os.Setenv("GOSH_WORKING_DIR", suite.workdir)
	goshDir := filepath.Join(suite.homedir, ".gosh")
	_ = os.MkdirAll(goshDir, 0755)
}

func (suite *ConfigTestSuite) TestInitializeConfig_NoConfig() {
	InitializeConfig()
	//r := suite.Require()
}

func (suite *ConfigTestSuite) TestInitializeConfig_BasicAuth_ConfigFile() {
	contents := []byte(`
Auth:
  Type: basic
  User: user1
  Pass: my-pass
`)
	r := suite.Require()
	err := os.WriteFile(filepath.Join(suite.homedir, ".gosh", "config.yml"), contents, 0644)
	if err != nil {
		r.Fail("unable to init test, cannot create ~/.gosh/config.yml file")
	}
	InitializeConfig()

	r.Equal(BasicAuth, Config.Auth.Type())
	auth := Config.Auth.(BasicAuthConfig)
	r.Equal("user1", auth.User)
	r.Equal("my-pass", auth.Pass)
}

func (suite *ConfigTestSuite) TestInitializeConfig_BasicAuth_Env() {
	_ = os.Setenv("GOSH_AUTH_TYPE", "basic")
	_ = os.Setenv("GOSH_AUTH_USER", "user1")
	_ = os.Setenv("GOSH_AUTH_PASS", "my-pass")

	InitializeConfig()
	r := suite.Require()
	r.Equal(BasicAuth, Config.Auth.Type())
	auth := Config.Auth.(BasicAuthConfig)
	r.Equal("user1", auth.User)
	r.Equal("my-pass", auth.Pass)
}

func (suite *ConfigTestSuite) TestInitializeOutputConfig_ConfigFile() {
	contents := []byte(`
Output:
  default_format: properties
  versions_key_suffix: version
`)
	r := suite.Require()
	err := os.WriteFile(filepath.Join(suite.homedir, ".gosh", "config.yml"), contents, 0644)
	if err != nil {
		r.Fail("unable to init test, cannot create ~/.gosh/config.yml file")
	}
	InitializeConfig()
	r.Equal("properties", Config.Output.DefaultFormat)
	r.Equal("version", Config.Output.VersionsKeySuffix)
	r.Equal("", Config.Output.ArtifactsKeySuffix)
}

func (suite *ConfigTestSuite) TestInitializeConfig_SshAuth_ConfigFile() {
	contents := []byte(`
Auth:
  Type: ssh
  Private_Key_File: private-key-file
  Private_Key_Pass: private-key-pass
`)
	r := suite.Require()
	err := os.WriteFile(filepath.Join(suite.homedir, ".gosh", "config.yml"), contents, 0644)
	if err != nil {
		r.Fail("unable to init test, cannot create ~/.gosh/config.yml file")
	}
	InitializeConfig()

	r.Equal(SshKey, Config.Auth.Type())
	auth := Config.Auth.(SshAuthConfig)
	r.Equal("private-key-file", auth.PrivateKeyFile)
	r.Equal("private-key-pass", auth.PrivateKeyPass)
}

func (suite *ConfigTestSuite) TestInitializeConfig_SshAuth_Env() {
	_ = os.Setenv("GOSH_AUTH_TYPE", "ssh")
	_ = os.Setenv("GOSH_AUTH_PRIVATE_KEY_FILE", "private-key-file")
	_ = os.Setenv("GOSH_AUTH_PRIVATE_KEY_PASS", "private-key-pass")

	InitializeConfig()
	r := suite.Require()
	r.Equal(SshKey, Config.Auth.Type())
	auth := Config.Auth.(SshAuthConfig)
	r.Equal("private-key-file", auth.PrivateKeyFile)
	r.Equal("private-key-pass", auth.PrivateKeyPass)
}

func (suite *ConfigTestSuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
