package gitops

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"gosh/util"
	"path/filepath"
	"testing"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type AppSuite struct {
	suite.Suite
	appGroup *AppGroup
	app      *App
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *AppSuite) SetupSuite() {
	TestsSetupWorkingDir(suite.Suite)
	CreateTestAppGroup(suite.Suite, "test")
	suite.appGroup = &AppGroup{Name: "test"}
	CreateTestApp(suite.Suite, "app1", "test")
	suite.app = &App{Name: "app1", group: suite.appGroup}

}

func (suite *AppSuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

func (suite *AppSuite) TestGetFilePath() {
	app := NewApp("app1", suite.appGroup)
	r := suite.Require()
	r.Equal(filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/test/app1.yml"), app.GetFilePath())
}

func (suite *AppSuite) TestReadInvalidStructReturnValidationErr() {
	app := &App{}
	err := app.Read()
	r := suite.Require()
	r.NotNil(err)
	r.Equal(ValidationErr, err)
}

func (suite *AppSuite) TestRead() {
	app := NewApp("app1", suite.appGroup)
	err := app.Read()
	r := suite.Require()
	r.Nil(err)
	r.Len(app.Properties, 1)
	r.Equal("app1", app.Properties["app_name"])
}

func (suite *AppSuite) TestCreate() {
	app := NewApp("app-create", suite.appGroup)
	app.Properties["groupId"] = "com/trendminer/group"
	app.Properties["groupId"] = "app-create-dist"
	app.Artifacts["test"] = "https://nexus.trendminer.net/repository/app-create-dist.zip"
	err := app.Create()
	r := suite.Require()
	r.Nil(err)
	app2 := NewApp("app-create", suite.appGroup)
	err = app2.Read()
	r.Nil(err)
	r.Equal(app.Name, app2.Name)
	r.Equal(app.Properties, app2.Properties)
	r.Equal(app.Artifacts, app2.Artifacts)
	g := NewAppGroup("test")
	_ = g.Read()
	found := false
	for _, v := range g.Apps {
		found = v.Name == app.Name
	}
	r.True(found, "app not found in group")
}

func (suite *AppSuite) TestCreateFromDefaultTemplate() {
	app := NewApp("app-create-default-templ", suite.appGroup)
	app.Properties["groupId"] = "com/trendminer/group"
	app.Properties["artifactId"] = "app-create-dist"
	err := app.CreateFromTemplate("")
	r := suite.Require()
	r.Nil(err)
	app2 := NewApp("app-create-default-templ", suite.appGroup)
	err = app2.Read()
	r.Nil(err)
	r.Equal("app-create-default-templ", app2.Name)
	r.Len(app2.Properties, 5)
	r.Len(app2.Artifacts, 2)
	r.Equal("[gosh:repo:docker]/app-create-default-templ:[gosh:version]", app2.Artifacts["docker"])
}

func (suite *AppSuite) TestFindApp() {
	CreateTestApp(suite.Suite, "app2", suite.appGroup.Name)
	app, err := FindApp("app2")
	r := suite.Require()
	r.Nil(err)
	r.NotNil(app)
	r.Equal("app2", app.Name)
}

func (suite *AppSuite) TestFindAppNonExistingReturnsErr() {
	_, err := FindApp("non-existing-app")
	r := suite.Require()
	r.NotNil(err)
	r.Equal(ResourceDoesNotExistErr, err)
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}
