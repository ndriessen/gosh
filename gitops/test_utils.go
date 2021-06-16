package gitops

import (
	"bytes"
	"github.com/Flaque/filet"
	"github.com/stretchr/testify/suite"
	"gosh/util"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const testAppGroupFileContents = `
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

const testStageFileContents = `
parameters:
  {{.Name}}:
    app1: 1.0.0
    app2: 2.0.0
    app3: 3.0.0
`

const testAppFileContents = `
parameters:
  {{.Name}}:
    app_name: {{.Name}}
`

const testReleaseFileContents = `
parameters:
  {{.Name}}:
    app1: 
      version: 1.0.0
    app2: 
      version: 2.0.0
    app3:
      version: 3.0.0
`

func init() {
	appGroupContentsTemplate, _ = template.New("appgroup").Parse(testAppGroupFileContents)
	stageContentsTemplate, _ = template.New("stage").Parse(testStageFileContents)
	appContentsTemplate, _ = template.New("app").Parse(testAppFileContents)
	releaseContentsTemplate, _ = template.New("release").Parse(testReleaseFileContents)
	testConfig = &util.GoshConfig{
		DeploymentRepository: util.DeploymentRepository{
			Url:               "",
			SshKey:            "",
			SshPrivateKeyPass: "",
		},
		ArtifactRepositories: nil,
	}
}

var appContentsTemplate, appGroupContentsTemplate, stageContentsTemplate, releaseContentsTemplate *template.Template
var testConfig *util.GoshConfig

func TestsSetupWorkingDir(suite suite.Suite) {
	dir := filet.TmpDir(suite.T(), "")
	util.Context.WorkingDir = dir
	p := filepath.Join(util.Context.WorkingDir, "inventory/classes/releases/stage")
	_ = os.MkdirAll(p, 0755)
	p = filepath.Join(util.Context.WorkingDir, "inventory/classes/releases/product")
	_ = os.MkdirAll(p, 0755)
	p = filepath.Join(util.Context.WorkingDir, "inventory/classes/releases/hotfix")
	_ = os.MkdirAll(p, 0755)
	p = filepath.Join(util.Context.WorkingDir, "inventory/classes/apps")
	_ = os.MkdirAll(p, 0755)
	p = filepath.Join(util.Context.WorkingDir, "inventory/classes/stages")
	_ = os.MkdirAll(p, 0755)
}

func CreateTestAppGroup(suite suite.Suite, name string) {
	if name == "" {
		name = "test"
	}
	p := filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/", name)
	_ = os.MkdirAll(p, 0755)
	f := filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/", name+".yml")
	var tpl bytes.Buffer
	if err := appGroupContentsTemplate.Execute(&tpl, &AppGroup{Name: name}); err == nil {
		filet.File(suite.T(), f, tpl.String())
	} else {
		log.Fatalln("Could not create test app group", err)
	}
}

func CreateTestStage(suite suite.Suite, name string) {
	if name == "" {
		name = "alpha"
	}
	f := filepath.Join(util.Context.WorkingDir, "inventory/classes/stages/", name+".yml")
	var tpl bytes.Buffer
	if err := stageContentsTemplate.Execute(&tpl, &Stage{Name: name}); err == nil {
		filet.File(suite.T(), f, tpl.String())
	} else {
		log.Fatalln("Could not create test stage", err)
	}
}

func CreateTestRelease(suite suite.Suite, name string, rType ReleaseType) {
	if name == "" {
		name = "my-release"
	}
	f := filepath.Join(util.Context.WorkingDir, "inventory/classes/releases/", rType.String(), name+".yml")
	var tpl bytes.Buffer
	if err := releaseContentsTemplate.Execute(&tpl, &Release{Name: name}); err == nil {
		filet.File(suite.T(), f, tpl.String())
	} else {
		log.Fatalln("Could not create test release", err)
	}
}

func CreateTestApp(suite suite.Suite, name string, group string) {
	if name == "" {
		name = "test-app"
	}
	f := filepath.Join(util.Context.WorkingDir, "inventory/classes/apps/", group, name+".yml")
	var tpl bytes.Buffer
	if err := appContentsTemplate.Execute(&tpl, NewApp(name, NewAppGroup("test"))); err == nil {
		filet.File(suite.T(), f, tpl.String())
	} else {
		log.Fatalln("Could not create test app", err)
	}
}
