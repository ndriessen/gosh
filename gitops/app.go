package gitops

import (
	"errors"
	"gosh/log"
	"gosh/util"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var (
	NoSuchArtifactErr = errors.New("no such artifact")
)

type App struct {
	Name       string
	Properties map[string]string
	Artifacts  map[string]string
	group      *AppGroup
}

func NewApp(name string, group *AppGroup) *App {
	return &App{Name: name, group: group, Properties: map[string]string{}, Artifacts: map[string]string{}}
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
	return app, err
}

func (app *App) String() string {
	if app == nil {
		return "[app: nil]"
	}
	return "[app: " + app.Name + "]"
}

func (app *App) Read() error {
	return Read(app)
}

func (app *App) Create() (err error) {
	if err = prepareCreate(app); err == nil {
		if err = createFromStruct(app); err == nil {
			return addAppToGroup(app)
		}
	}
	return
}
func (app *App) CreateFromTemplate(templateName string) (err error) {
	if err = prepareCreate(app); err == nil {
		if err = createFromTemplate(app, templateName); err == nil {
			return addAppToGroup(app)
		}
	}
	return
}

func addAppToGroup(app *App) error {
	if err := app.group.Read(); err == nil {
		app.group.Apps = append(app.group.Apps, app)
		return app.group.Update()
	} else {
		return log.Errf(err, "could not add app %s to group", app.Name)
	}
}

func prepareCreate(app *App) error {
	log.Tracef("Create app with input: %+v", app)
	if app == nil || !app.isValid() {
		return log.Err(ValidationErr, "Invalid app struct, use NewApp() to create one")
	}
	if app.Exists() {
		return log.Errf(ResourceAlreadyExistsErr, "The app '%s' already exists", app.Name)
	}
	if existing, _ := FindApp(app.Name); existing != nil {
		return log.Errf(ResourceAlreadyExistsErr, "The app '%s' already exists in another group, app names must be unique", app.Name)
	}
	if !app.group.Exists() {
		log.Debugf("Creating group %s for app %s", app.group.Name, app.Name)
		if err := app.group.Create(); err != nil {
			return log.Errf(err, "Error creating group %s for app %s", app.group.Name, app.Name)
		}
	}
	return nil
}

func createFromTemplate(app *App, templateName string) error {
	if t, err := NewAppTemplate(templateName); err == nil {
		if f, err := t.Render(app); err == nil {
			if err = ioutil.WriteFile(app.GetFilePath(), []byte(f), 0644); err == nil {
				return nil
			} else {
				return log.Errf(err, "error writing app file")
			}
		} else {
			return log.Errf(err, "error rendering app template %s", templateName)
		}
	} else {
		return log.Errf(err, "error loading app template %s", templateName)
	}
}

func createFromStruct(app *App) error {
	f := app.mapToKapitanFile()
	if err := WriteKapitanFile(app.GetFilePath(), f); err == nil {
		log.Infof("Created app '%s", app.Name)
		return nil
	} else {
		return log.Errf(err, "Error writing app group file '%s'", app.GetFilePath())
	}
}

func (app *App) mapToKapitanFile() *kapitanFile {
	log.Tracef("Mapping app %s to kapitan file: %+v", app.Name, app)
	f := newKapitanFile()
	props := map[string]interface{}{}
	for key, value := range app.Properties {
		props[key] = value
	}
	props["artifacts"] = app.Artifacts
	f.Parameters[app.Name] = props
	log.Tracef("Mapped app %s to kapitan file, result: %+v", app.Name, f)
	return f
}

func (app *App) mapFromKapitanFile(f *kapitanFile) {
	log.Tracef("Mapping app %s from kapitan file %+v", app.Name, f)
	app.Properties = make(map[string]string, 0)
	app.Artifacts = make(map[string]string, 0)
	if properties, exists := f.Parameters[app.Name]; exists {
		for key, value := range properties.(map[interface{}]interface{}) {
			if key == "artifacts" {
				for t, u := range properties.(map[interface{}]interface{})[key].(map[interface{}]interface{}) {
					app.Artifacts[t.(string)] = u.(string)
				}
			} else {
				switch value.(type) {
				case string:
					app.Properties[key.(string)] = value.(string)
				default:
					log.Warnf("app definition '%s' has nested key '%s', not mapping", app.Name, key.(string))
				}
			}
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

func (app *App) GetArtifact(list VersionsList, version string, artifactType string) (string, error) {
	if app.Artifacts != nil {
		if value, exists := app.Artifacts[artifactType]; exists {
			if replacements, err := buildReplacementMap(list, version); err == nil {
				for find, replace := range replacements {
					value = strings.ReplaceAll(value, find, replace)
				}
				return value, nil
			} else {
				return "", log.Errf(errors.New("error replacing gosh placeholder in artifact"), "Could not generate artifact list")
			}
		}
	}
	return "", NoSuchArtifactErr
}

func buildReplacementMap(list VersionsList, version string) (map[string]string, error) {
	result := map[string]string{}
	for t, urls := range util.Config.ArtifactRepositories {
		if v, exists := urls[strings.ToLower(list.getResourceName())]; exists {
			result["[gosh:repo:"+t+"]"] = v
		} else if v, exists = urls["default"]; exists {
			result["[gosh:repo:"+t+"]"] = v
		} else {
			return nil, log.Errf(errors.New("No default repository set for type "+t), "could not generate artifacts")
		}
	}
	result["[gosh:version]"] = version
	return result, nil
}
