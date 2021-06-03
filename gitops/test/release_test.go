package gitops_test

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"gosh/gitops"
	"gosh/util"
	"os"
	"path/filepath"
	"testing"
)

const testReleaseFileContents = `
parameters:
  R2021.R1:
    app1:
      version: 1.0.0
    app2:
      version: 2.0.0
    app3:
      version: 3.0.0
`

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type ReleaseSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *ReleaseSuite) SetupSuite() {
	dir := filet.TmpDir(suite.T(), "")
	util.Context.WorkingDir = dir
	p := filepath.Join(util.Context.WorkingDir, "inventory/classes/releases")
	_ = os.MkdirAll(p, 0755)
	p = filepath.Join(util.Context.WorkingDir, "inventory/classes/releases/product")
	_ = os.MkdirAll(p, 0755)
	p = filepath.Join(util.Context.WorkingDir, "inventory/classes/releases/hotfix")
	_ = os.MkdirAll(p, 0755)
}

func (suite *ReleaseSuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

func (suite *ReleaseSuite) TestGetFilePath() {
	release := gitops.NewRelease("R2021.R1", gitops.ProductRelease)
	r := suite.Require()
	r.Equal(filepath.Join(util.Context.WorkingDir, "inventory/classes/releases/product/R2021.R1.yml"), release.GetFilePath())
}

func (suite *ReleaseSuite) TestReadInvalidStructReturnValidationErr() {
	release := &gitops.Release{}
	err := release.Read()
	r := suite.Require()
	r.NotNil(err)
	r.Equal(gitops.ValidationErr, err)
}

func (suite *ReleaseSuite) TestRead() {
	f := filepath.Join(util.Context.WorkingDir, "inventory/classes/releases/product/R2021.R1.yml")
	filet.File(suite.T(), f, testReleaseFileContents)
	release := gitops.NewRelease("R2021.R1", gitops.ProductRelease)
	err := release.Read()
	r := suite.Require()
	r.Nil(err)
	r.Len(release.Versions, 3)
	r.Equal("R2021.R1", release.Name)
	r.Equal("1.0.0", release.Versions["app1"])
	r.Equal("2.0.0", release.Versions["app2"])
	r.Equal("3.0.0", release.Versions["app3"])
}

func TestReleaseTestSuite(t *testing.T) {
	suite.Run(t, new(ReleaseSuite))
}
