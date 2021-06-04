package gitops

import (
	"bytes"
	"errors"
	"gosh/log"
	"gosh/util"
	"os"
	"path/filepath"
	"text/template"
)

const (
	defaultAppTemplateContents = `parameters:
  {{.Name}}:
    app_name: {{.Name}}
    version: latest
    ref: master
    groupId: {{.Properties.groupId}}
    artifactId: {{.Properties.artifactId}}
    artifacts:
      maven: "[gosh:repo:maven]/{{.Properties.groupId}}/{{.Properties.artifactId}}/[gosh:version]/{{.Properties.artifactId}}-[gosh:version].zip"
      docker: "[gosh:repo:docker]/ {{- .Name -}} :[gosh:version]"
`
)

var (
	DefaultAppTemplate  *AppTemplate
	TemplateNotFoundErr = errors.New("template does not exist")
	TemplateRenderErr   = errors.New("error rendering template")
)

func init() {
	if t, err := template.New("default").Parse(defaultAppTemplateContents); err == nil {
		DefaultAppTemplate = &AppTemplate{
			Name:     "Default",
			template: t,
		}
	} else {
		_ = log.Errf(err, "Could not load default app template")
	}
}

type AppTemplateAction interface {
	Render(app *App) (string, error)
}

type AppTemplate struct {
	Name     string
	template *template.Template
}

func NewAppTemplate(templateName string) (*AppTemplate, error) {
	if templateName == "" {
		return DefaultAppTemplate, nil
	}
	templateFile := filepath.Join(util.Context.WorkingDir, ".gosh", "templates", templateName+".yml")
	if info, err := os.Stat(templateFile); err != nil && !info.IsDir() {
		if t, err := template.New("default").Parse(defaultAppTemplateContents); err == nil {
			return &AppTemplate{
				Name:     templateName,
				template: t,
			}, nil
		} else {
			return nil, log.Errf(err, "Could not load app template %s", templateName)
		}
	} else {
		return nil, log.Errf(TemplateNotFoundErr, "Could not find template %s", templateFile)
	}
}

func (t *AppTemplate) Render(app *App) (string, error) {
	result := new(bytes.Buffer)
	if err := t.template.Execute(result, app); err == nil {
		return result.String(), nil
	}
	return "", log.Errf(TemplateRenderErr, "Could not render template %s with data %+v", t.Name, app)
}
