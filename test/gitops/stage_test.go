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

type StageSuite struct {
	suite.Suite
	stage *gitops.Stage
}

func (suite *StageSuite) SetupSuite() {
	test.SetupWorkingDir(suite.Suite)
	test.CreateTestStage(suite.Suite, "alpha")
	suite.stage = gitops.NewStage("alpha")
	test.CreateTestAppGroup(suite.Suite, "test")
	test.CreateTestApp(suite.Suite, "test-app", "test")
	test.CreateTestRelease(suite.Suite, "alpha", gitops.StageRelease)
}

func (suite *StageSuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

func (suite *StageSuite) TestGetFilePath() {
	stage := gitops.NewStage("alpha")
	r := suite.Require()
	r.Equal(filepath.Join(util.Context.WorkingDir, "inventory/classes/stages/alpha.yml"), stage.GetFilePath())
}

func (suite *StageSuite) TestReadInvalidStructReturnValidationErr() {
	app := &gitops.App{}
	err := app.Read()
	r := suite.Require()
	r.NotNil(err)
	r.Equal(gitops.ValidationErr, err)
}

func (suite *StageSuite) TestRead() {
	stage := gitops.NewStage("alpha")
	err := stage.Read()
	r := suite.Require()
	r.Nil(err)
	r.Len(stage.Versions, 3)
	r.Equal("alpha", stage.Name)
	r.Equal("1.0.0", stage.Versions["app1"])
	r.Equal("2.0.0", stage.Versions["app2"])
	r.Equal("3.0.0", stage.Versions["app3"])
}

func (suite *StageSuite) TestCreate() {
	stage := gitops.NewStage("mystage")
	stage.Versions["app1"] = "1.0.0"
	err := stage.Create()
	r := suite.Require()
	r.Nil(err)
	r.Equal("mystage", stage.Name)
	r.Len(stage.Versions, 1)
	r.Equal("1.0.0", stage.Versions["app1"])

	release := gitops.NewRelease("mystage", gitops.StageRelease)
	r.True(release.Exists())
	err = release.Read()
	r.Nil(err)
	r.Equal("mystage", release.Name)
	r.Len(release.Versions, 1)
	r.Equal("1.0.0", release.Versions["app1"])
}

func (suite *StageSuite) TestUpdate() {
	r := suite.Require()
	stage := gitops.NewStage("alpha")
	err := stage.Read()
	r.Nil(err)
	release := gitops.NewRelease("alpha", gitops.StageRelease)
	err = stage.Read()
	r.Nil(err)

	err = stage.UpdateVersion("test-app", "1.2.3")
	r.Nil(err)
	r.Len(stage.Versions, 4)
	r.Equal("alpha", stage.Name)
	r.Equal("1.0.0", stage.Versions["app1"])
	r.Equal("2.0.0", stage.Versions["app2"])
	r.Equal("3.0.0", stage.Versions["app3"])
	r.Equal("1.2.3", stage.Versions["test-app"])

	err = release.Read()
	r.Nil(err)
	r.Len(release.Versions, 4)
	r.Equal("alpha", release.Name)
	r.Equal("1.0.0", release.Versions["app1"])
	r.Equal("2.0.0", release.Versions["app2"])
	r.Equal("3.0.0", release.Versions["app3"])
	r.Equal("1.2.3", release.Versions["test-app"])
}

func (suite *StageSuite) TestUpdate_NotReadFirst() {
	err := suite.stage.Update()
	r := suite.Require()
	r.NotNil(err)
	r.Equal(gitops.ResourceUpdatedWithoutReadingErr, err)
}

func TestStageTestSuite(t *testing.T) {
	suite.Run(t, new(StageSuite))
}
