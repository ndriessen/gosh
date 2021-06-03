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

const testAppGroupFileContents = `
classes:
  - apps.test.app1
  - apps.test.app2
parameters:
  test:
    prop1: value1
    prop2:
      nested1: value2
  other:
    props:
    - array

`

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type AppGroupSuite struct {
	suite.Suite
	TestAppGroupName string
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *AppGroupSuite) SetupSuite() {
	dir := filet.TmpDir(suite.T(), "")
	util.Context.WorkingDir = dir
	p := filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/test")
	_ = os.MkdirAll(p, 0755)
}

func (suite *AppGroupSuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

func (suite *AppGroupSuite) TestNewAppGroup() {
	group := gitops.NewAppGroup("test")
	r := suite.Require()
	r.Equal("test", group.Name)
}

func (suite *AppGroupSuite) TestGetFilePath() {
	group := gitops.NewAppGroup("test")
	r := suite.Require()
	r.Equal(filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/test.yml"), group.GetFilePath())
}

func (suite *AppGroupSuite) TestGetFolderPath() {
	group := gitops.NewAppGroup("test")
	r := suite.Require()
	r.Equal(filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/test"), group.GetFolderPath())
}

func (suite *AppGroupSuite) TestCreate() {
	group := gitops.NewAppGroup("my-group", &gitops.App{Name: "app1"}, &gitops.App{Name: "app2"})
	err := group.Create()
	f, _ := gitops.ReadKapitanFile(group.GetFilePath())
	r := suite.Require()
	r.Nil(err)
	r.Len(f.Classes, 2)
	r.Equal("apps.my-group.app1", f.Classes[0])
	r.Equal("apps.my-group.app2", f.Classes[1])
}

func (suite *AppGroupSuite) TestRead() {
	f := filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/test.yml")
	filet.File(suite.T(), f, testAppGroupFileContents)
	group := gitops.NewAppGroup("test")
	err := group.Read()
	r := suite.Require()
	r.Nil(err)
	r.Equal("test", group.Name)
	r.Len(group.Apps, 2)
	r.Equal("app1", group.Apps[0].Name)
	r.Equal("app2", group.Apps[1].Name)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAppGroupTestSuite(t *testing.T) {
	suite.Run(t, new(AppGroupSuite))
}
