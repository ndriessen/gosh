package git

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"gosh/gitops"
	"testing"
)

type DeploymentRepositorySuite struct {
	suite.Suite
}

func (suite *DeploymentRepositorySuite) SetupSuite() {
	gitops.TestsSetupWorkingDir(suite.Suite)
}

func (suite *DeploymentRepositorySuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

//todo: we need mocking to test this...

func TestDeploymentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DeploymentRepositorySuite))
}
