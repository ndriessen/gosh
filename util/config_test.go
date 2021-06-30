package util

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) SetupSuite() {
}

func (suite *ConfigTestSuite) TestInitializeConfig_NoConfig() {
	os.Setenv("HOME", filet.TmpDir(suite.T(), ""))
	InitializeConfig()
	//r := suite.Require()
}

func (suite *ConfigTestSuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
