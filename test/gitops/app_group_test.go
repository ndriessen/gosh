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

type AppGroupSuite struct {
	suite.Suite
	TestAppGroupName string
}

func (suite *AppGroupSuite) SetupSuite() {
	test.SetupWorkingDir(suite.Suite)
	test.CreateTestAppGroup(suite.Suite, "test")
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
