package gitops_test

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"gosh/gitops"
	"gosh/test"
	"testing"
)

type VersionsListSuite struct {
	suite.Suite
	testStage *gitops.Stage
}

func (suite *VersionsListSuite) SetupSuite() {
	test.SetupWorkingDir(suite.Suite)
	test.CreateTestAppGroup(suite.Suite, "test")
	versions := map[string]string{}
	versions["app1"] = "1.0.0"
	versions["app2"] = "2.0.0"
	versions["app3"] = "3.0.0"

	suite.testStage = gitops.NewStage("alpha")
	suite.testStage.Versions = versions

}

func (suite *VersionsListSuite) TestGetVersionsNoFilter() {
	v := gitops.GetVersions(suite.testStage, "", "")
	r := suite.Require()
	r.NotNil(v)
	r.Len(v, 3)
}

func (suite *VersionsListSuite) TestGetVersionsBothFilters() {
	v := gitops.GetVersions(suite.testStage, "test", "app1")
	r := suite.Require()
	r.NotNil(v)
	r.Len(v, 1)
	a, e := v["app1"]
	r.True(e)
	r.Equal("1.0.0", a)
}

func (suite *VersionsListSuite) TestGetVersionsAppFilter() {
	v := gitops.GetVersions(suite.testStage, "", "app1")
	r := suite.Require()
	r.NotNil(v)
	r.Len(v, 1)
	a, e := v["app1"]
	r.True(e)
	r.Equal("1.0.0", a)
}

func (suite *VersionsListSuite) TestGetVersionsGroupFilter() {
	v := gitops.GetVersions(suite.testStage, "test", "")
	r := suite.Require()
	r.NotNil(v)
	r.Len(v, 2)
	a, e := v["app1"]
	r.True(e)
	r.Equal("1.0.0", a)
	a, e = v["app2"]
	r.True(e)
	r.Equal("2.0.0", a)
}

func (suite *VersionsListSuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

func TestVersionsListTestSuite(t *testing.T) {
	suite.Run(t, new(VersionsListSuite))
}
