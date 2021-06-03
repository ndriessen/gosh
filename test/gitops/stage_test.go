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

const testStageFileContents = `
parameters:
  alpha:
    app1: 1.0.0
    app2: 2.0.0
    app3: 3.0.0
`

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type StageSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *StageSuite) SetupSuite() {
	dir := filet.TmpDir(suite.T(), "")
	util.Context.WorkingDir = dir
	p := filepath.Join(util.Context.WorkingDir, "inventory/classes/stages")
	_ = os.MkdirAll(p, 0755)
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
	f := filepath.Join(util.Context.WorkingDir, "inventory/classes/stages/alpha.yml")
	filet.File(suite.T(), f, testStageFileContents)
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

func TestStageTestSuite(t *testing.T) {
	suite.Run(t, new(StageSuite))
}
