package test

import (
	"bytes"
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"gosh/gitops"
	"gosh/util"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const TestAppGroupFileContents = `
classes:
  - apps.{{.Name}}.app1
  - apps.{{.Name}}.app2
parameters:
  {{.Name}}:
    prop1: value1
    prop2:
      nested1: value2
  other:
    props:
    - array

`

func init() {
	appGroupContentsTemplate, _ = template.New("appgroup").Parse(TestAppGroupFileContents)
}

var appGroupContentsTemplate *template.Template

func SetupWorkingDir(suite suite.Suite) {
	dir := filet.TmpDir(suite.T(), "")
	util.Context.WorkingDir = dir
}

func CreateTestAppGroup(suite suite.Suite, name string) {
	if name == "" {
		name = "test"
	}
	p := filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/test")
	_ = os.MkdirAll(p, 0755)
	f := filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/test.yml")
	var tpl bytes.Buffer
	if err := appGroupContentsTemplate.Execute(&tpl, &gitops.AppGroup{Name: name}); err == nil {
		filet.File(suite.T(), f, tpl.String())
	} else {
		log.Fatalln("Could not create test app group", err)
	}
}
