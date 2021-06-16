package gitops

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"gosh/util"
	"path/filepath"
	"testing"
)

type StageSuite struct {
	suite.Suite
	stage *Stage
}

func (suite *StageSuite) SetupSuite() {
	TestsSetupWorkingDir(suite.Suite)
	CreateTestStage(suite.Suite, "alpha")
	suite.stage = NewStage("alpha")
	CreateTestAppGroup(suite.Suite, "test")
	CreateTestApp(suite.Suite, "test-app", "test")
	CreateTestRelease(suite.Suite, "alpha", StageRelease)
}

func (suite *StageSuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

func (suite *StageSuite) TestGetFilePath() {
	stage := NewStage("alpha")
	r := suite.Require()
	r.Equal(filepath.Join(util.Context.WorkingDir, "inventory/classes/stages/alpha.yml"), stage.GetFilePath())
}

func (suite *StageSuite) TestReadInvalidStructReturnValidationErr() {
	app := &App{}
	err := app.Read()
	r := suite.Require()
	r.NotNil(err)
	r.Equal(ValidationErr, err)
}

func (suite *StageSuite) TestRead() {
	stage := NewStage("alpha")
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
	stage := NewStage("mystage")
	stage.Versions["app1"] = "1.0.0"
	err := stage.Create()
	r := suite.Require()
	r.Nil(err)
	r.Equal("mystage", stage.Name)
	r.Len(stage.Versions, 1)
	r.Equal("1.0.0", stage.Versions["app1"])

	release := NewRelease("mystage", StageRelease)
	r.True(release.Exists())
	err = release.Read()
	r.Nil(err)
	r.Equal("mystage", release.Name)
	r.Len(release.Versions, 1)
	r.Equal("1.0.0", release.Versions["app1"])
}

func (suite *StageSuite) TestUpdate() {
	r := suite.Require()
	stage := NewStage("alpha")
	err := stage.Read()
	r.Nil(err)
	release := NewRelease("alpha", StageRelease)
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
	r.Equal(ResourceUpdatedWithoutReadingErr, err)
}

func TestStageTestSuite(t *testing.T) {
	suite.Run(t, new(StageSuite))
}
