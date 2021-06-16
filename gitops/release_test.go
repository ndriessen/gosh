package gitops

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"gosh/util"
	"path/filepath"
	"testing"
)

type ReleaseSuite struct {
	suite.Suite
}

func (suite *ReleaseSuite) SetupSuite() {
	TestsSetupWorkingDir(suite.Suite)
	CreateTestRelease(suite.Suite, "test-release", ProductRelease)
}

func (suite *ReleaseSuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

func (suite *ReleaseSuite) TestGetFilePath() {
	release := NewRelease("R2021.R1", ProductRelease)
	r := suite.Require()
	r.Equal(filepath.Join(util.Context.WorkingDir, "inventory/classes/releases/product/R2021.R1.yml"), release.GetFilePath())
}

func (suite *ReleaseSuite) TestReadInvalidStructReturnValidationErr() {
	release := &Release{}
	err := release.Read()
	r := suite.Require()
	r.NotNil(err)
	r.Equal(ValidationErr, err)
}

func (suite *ReleaseSuite) TestRead() {
	release := NewRelease("test-release", ProductRelease)
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
	_, err := NewReleaseFromFullName("invalid")
	r := suite.Require()
	r.NotNil(err)
	r.Equal(InvalidFullReleaseNameErr, err)
}

func (suite *ReleaseSuite) TestNewReleaseFromFullNameInvalidType() {
	_, err := NewReleaseFromFullName("invalid/release")
	r := suite.Require()
	r.NotNil(err)
	r.Equal(UnsupportedReleaseTypeErr, err)
}

func (suite *ReleaseSuite) TestNewReleaseFromFullNameStageReleaseNotSupported() {
	_, err := NewReleaseFromFullName("stage/release")
	r := suite.Require()
	r.NotNil(err)
	r.Equal(InvalidFullReleaseNameErr, err)
}

func (suite *ReleaseSuite) TestNewReleaseFromFullName() {
	release, err := NewReleaseFromFullName("product/release")
	r := suite.Require()
	r.Nil(err)
	r.Equal(ProductRelease, release.Type)
	r.Equal("release", release.Name)
}

func (suite *ReleaseSuite) TestCreate() {
	release := NewRelease("R1", StageRelease)
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
	release := NewRelease("test-release", ProductRelease)
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
