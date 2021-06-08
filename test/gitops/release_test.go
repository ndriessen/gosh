package gitops_test

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"gosh/gitops"
	"gosh/test"
	"gosh/util"
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
	test.SetupWorkingDir(suite.Suite)
	test.CreateTestRelease(suite.Suite, "test-release", gitops.ProductRelease)
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
	release := gitops.NewRelease("test-release", gitops.ProductRelease)
	err := release.Read()
	r := suite.Require()
	r.Nil(err)
	r.Len(release.Versions, 3)
	r.Equal("test-release", release.Name)
	r.Equal("1.0.0", release.Versions["app1"])
	r.Equal("2.0.0", release.Versions["app2"])
	r.Equal("3.0.0", release.Versions["app3"])
}

func (suite *ReleaseSuite) TestNewReleaseFromFullNameInvalidName() {
	_, err := gitops.NewReleaseFromFullName("invalid")
	r := suite.Require()
	r.NotNil(err)
	r.Equal(gitops.InvalidFullReleaseNameErr, err)
}

func (suite *ReleaseSuite) TestNewReleaseFromFullNameInvalidType() {
	_, err := gitops.NewReleaseFromFullName("invalid/release")
	r := suite.Require()
	r.NotNil(err)
	r.Equal(gitops.UnsupportedReleaseTypeErr, err)
}

func (suite *ReleaseSuite) TestNewReleaseFromFullNameStageReleaseNotSupported() {
	_, err := gitops.NewReleaseFromFullName("stage/release")
	r := suite.Require()
	r.NotNil(err)
	r.Equal(gitops.InvalidFullReleaseNameErr, err)
}

func (suite *ReleaseSuite) TestNewReleaseFromFullName() {
	release, err := gitops.NewReleaseFromFullName("product/release")
	r := suite.Require()
	r.Nil(err)
	r.Equal(gitops.ProductRelease, release.Type)
	r.Equal("release", release.Name)
}

func (suite *ReleaseSuite) TestCreate() {
	release := gitops.NewRelease("R1", gitops.StageRelease)
	release.Versions["app1"] = "1.0.0"
	release.Versions["app2"] = "2.0.0"
	err := release.Create()
	r := suite.Require()
	r.Nil(err)

	err = release.Read()
	r.Nil(err)
	r.Equal("R1", release.Name)
	r.Len(release.Versions, 2)
	r.Equal("1.0.0", release.Versions["app1"])
	r.Equal("2.0.0", release.Versions["app2"])
}

func (suite *ReleaseSuite) TestUpdate() {
	release := gitops.NewRelease("test-release", gitops.ProductRelease)
	r := suite.Require()

	err := release.Read()
	r.Nil(err)
	release.Versions["my-app"] = "4.0.0"
	err = release.Update()
	r.Nil(err)
	r.Equal("test-release", release.Name)
	r.Len(release.Versions, 4)
	r.Equal("1.0.0", release.Versions["app1"])
	r.Equal("2.0.0", release.Versions["app2"])
	r.Equal("3.0.0", release.Versions["app3"])
	r.Equal("4.0.0", release.Versions["my-app"])
}

func TestReleaseTestSuite(t *testing.T) {
	suite.Run(t, new(ReleaseSuite))
}
