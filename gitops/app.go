package gitops

import (
	"gosh/log"
	"gosh/util"
	"io/fs"
	"path/filepath"
	"strings"
)

type App struct {
	Name       string
	Properties map[string]string
	group      *AppGroup
}

func NewApp(name string, group *AppGroup) *App {
	return &App{Name: name, group: group, Properties: map[string]string{}}
}

func FindApp(name string) (app *App, err error) {
	err = filepath.WalkDir(filepath.Join(util.Context.WorkingDir, appGroupPath), func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return log.Errf(err, "error searching app %s in path %s", name, path)
		}
		if !info.IsDir() {
			if strings.HasSuffix(path, name+kapitanFileExt) {
				dir, _ := filepath.Split(path)
				groupName := filepath.Base(dir)
				app = NewApp(name, NewAppGroup(groupName))
				return filepath.SkipDir
			}
		}
		return nil
	})
	if app == nil && err == nil {
		err = ResourceDoesNotExistErr
	}
	return app, log.Errf(err, "could not find app %s", name)
}

func (app *App) Read() error {
	return Read(app)
}

func (app *App) mapToKapitanFile() *kapitanFile {
	log.Tracef("Mapping app %s to kapitan file: %+v", app.Name, app)
	f := &kapitanFile{}
	props := f.Parameters[app.Name].(map[string]interface{})
	for key, value := range app.Properties {
		props[key] = value
	}
	log.Tracef("Mapped app %s to kapitan file, result: %+v", app.Name, f)
	return f
}

func (app *App) mapFromKapitanFile(f *kapitanFile) {
	log.Tracef("Mapping app %s from kapitan file %+v", app.Name, f)
	app.Properties = make(map[string]string, 0)
	if properties, exists := f.Parameters[app.Name]; exists {
		for key, value := range properties.(map[interface{}]interface{}) {
			app.Properties[key.(string)] = value.(string)
		}
	}
	log.Tracef("Mapped app %s from kapitan file, result: %+v", app.Name, app)
}

func (app *App) Exists() bool {
	return Exists(app)
}

func (app *App) GetFilePath() string {
	return filepath.Join(util.Context.WorkingDir, appGroupPath, app.group.Name, app.Name+kapitanFileExt)
}

func (app *App) isValid() bool {
	return strings.TrimSpace(app.Name) != "" && app.group != nil && app.group.isValid()
}

func (app *App) getResourceType() string {
	return "app"
}

func (app *App) getResourceName() string {
	return app.Name
}
