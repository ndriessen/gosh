package list_test

import (
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"gosh/list"
	"gosh/test"
	"testing"
)

type OutputFormatSuite struct {
	suite.Suite
	versions map[string]string
}

func (suite *OutputFormatSuite) SetupSuite() {
	test.SetupWorkingDir(suite.Suite)
	//test.CreateTestAppGroup(suite.Suite, "test")
	suite.versions = map[string]string{}
	suite.versions["app1"] = "1.0.0"
	suite.versions["app2"] = "2.0.0"
	suite.versions["app3"] = "3.0.0"
}

func (suite *OutputFormatSuite) TestRenderUnsupportedOutputFormat() {
	output, err := list.Render("unsupported", suite.versions)
	r := suite.Require()
	r.Empty(output)
	r.NotNil(err)
	r.Equal(list.UnsupportedOutputFormatErr, err)
}

func (suite *OutputFormatSuite) TestRenderYaml() {
	output, err := list.Render("yaml", suite.versions)
	r := suite.Require()
	r.NotEmpty(output)
	r.Nil(err)
	expected := `app1: 1.0.0
app2: 2.0.0
app3: 3.0.0
`
	r.Equal(expected, output)
}

func (suite *OutputFormatSuite) TestRenderProperties() {
	output, err := list.Render("properties", suite.versions)
	r := suite.Require()
	r.NotEmpty(output)
	r.Nil(err)
	expected := `app1=1.0.0
app2=2.0.0
app3=3.0.0
`
	r.Equal(expected, output)
}

func (suite *OutputFormatSuite) TearDownSuite() {
	filet.CleanUp(suite.T())
}

func TestOutputFormatTestSuite(t *testing.T) {
	suite.Run(t, new(OutputFormatSuite))
}
