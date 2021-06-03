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

const testAppFileContents = `
classes:
  - app
parameters:
  app1:
    app_name: app1
    image_name: ${tm:docker_repo}/${app1:app_name}
    image_tag: latest
    chart_ref: master
    type: PLATFORM
    groupId: com.trendminer.assets
    artifactId: app1
    artifactType: zip
    artifactClassifier: dist
  kapitan:
    vars:
      app_name: ${app1:app_name}
    dependencies:
      - type: git
        output_path: components/charts/${app1:app_name}
        source: ${tm:git:gitlab}/trendminer-platform/${app1:app_name}.git
        subdir: ${app1:app_name}-docker/src/main/resources/helm/${app1:app_name}
        ref: ${app1:chart_ref} #tag, commit, branch
    compile:
      - output_path: .
        input_type: helm
        input_paths:
          - components/charts/${app1:app_name}
        helm_values:
          image_name: ${app1:image_name}
          image_tag: ${app1:image_tag}
        helm_params:
          namespace: ${target_name}
          release_name: ${app1:app_name}

`

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type AppSuite struct {
	suite.Suite
	appGroup *gitops.AppGroup
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *AppSuite) SetupSuite() {
	dir := filet.TmpDir(suite.T(), "")
	util.Context.WorkingDir = dir
	suite.appGroup = &gitops.AppGroup{Name: "test"}
	p := filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/test")
	_ = os.MkdirAll(p, 0755)
}

func (suite *AppSuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

func (suite *AppSuite) TestGetFilePath() {
	app := gitops.NewApp("app1", suite.appGroup)
	r := suite.Require()
	r.Equal(filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/test/app1.yml"), app.GetFilePath())
}

func (suite *AppSuite) TestReadInvalidStructReturnValidationErr() {
	app := &gitops.App{}
	err := app.Read()
	r := suite.Require()
	r.NotNil(err)
	r.Equal(gitops.ValidationErr, err)
}

func (suite *AppSuite) TestRead() {
	f := filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/test/app1.yml")
	filet.File(suite.T(), f, testAppFileContents)
	app := gitops.NewApp("app1", suite.appGroup)
	err := app.Read()
	r := suite.Require()
	r.Nil(err)
	r.Len(app.Properties, 9)
	r.Equal("app1", app.Properties["artifactId"])
}

func (suite *AppSuite) TestFindApp() {
	f := filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/test/app2.yml")
	filet.File(suite.T(), f, testAppFileContents)
	app, err := gitops.FindApp("app2")
	r := suite.Require()
	r.Nil(err)
	r.NotNil(app)
	r.Equal("app2", app.Name)
}

func (suite *AppSuite) TestFindAppNonExistingReturnsErr() {
	_, err := gitops.FindApp("non-existing-app")
	r := suite.Require()
	r.NotNil(err)
	r.Equal(gitops.ResourceDoesNotExistErr, err)
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}
