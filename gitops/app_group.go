package gitops

import (
	"gosh/log"
	"gosh/util"
	"os"
	"path/filepath"
	"strings"
)

const (
	appGroupPath = kapitanClassesPath + "apps"
	appPrefix    = "apps."
)

type AppGroup struct {
	Name  string
	Apps  []*App
	_read bool
}

func (group *AppGroup) initialized() bool {
	return group._read
}

func (group *AppGroup) setInitialized() {
	group._read = true
}

// NewAppGroup /** Constructor for an AppGroup type
func NewAppGroup(name string, apps ...*App) *AppGroup {
	if apps == nil {
		apps = []*App{}
	}
	return &AppGroup{Name: name, Apps: apps}
}

func (group *AppGroup) Create() error {
	log.Tracef("Create app group with input: %+v", group)
	if err := os.MkdirAll(group.GetFolderPath(), 0755); err == nil {
		return create(group)
	} else {
		return log.Errf(err, "Error creating app group folder '%s'", group.GetFolderPath())
	}
}

func (group *AppGroup) mapToKapitanFile() *kapitanFile {
	log.Tracef("Mapping group %s to kapitan file: %+v", group.Name, group)
	f := newKapitanFile()
	for _, app := range group.Apps {
		f.Classes = append(f.Classes, appPrefix+group.Name+"."+app.Name)
	}
	log.Tracef("Mapped group %s to kapitan file, result: %+v", group.Name, f)
	return f
}

func (group *AppGroup) mapFromKapitanFile(f *kapitanFile) {
	log.Tracef("Mapping group %s from kapitan file %+v", group.Name, f)
	for _, class := range f.Classes {
		app := &App{
			Name:  strings.TrimPrefix(class, appPrefix+group.Name+"."),
			group: group,
		}
		group.Apps = append(group.Apps, app)
	}
	log.Tracef("Mapped group %s from kapitan file, result: %+v", group.Name, group)
}

func (group *AppGroup) Read() error {
	return read(group)
}

func (group *AppGroup) Update() error {
	return update(group)
}

func (group *AppGroup) Delete() error {
	return nil
}

func (group *AppGroup) Exists() bool {
	if f, err := os.Stat(group.GetFolderPath()); err == nil && f.IsDir() {
		if f, err = os.Stat(group.GetFilePath()); err == nil && !f.IsDir() {
			return true
		}
	}
	return false
}

func (group *AppGroup) GetFolderPath() string {
	return filepath.Join(util.Context.WorkingDir, appGroupPath, group.Name)
}

func (group *AppGroup) GetFilePath() string {
	return filepath.Join(util.Context.WorkingDir, appGroupPath, group.Name+kapitanFileExt)
}

func (group *AppGroup) isValid() bool {
	return strings.TrimSpace(group.Name) != ""
}

func (group *AppGroup) getResourceType() string {
	return "app group"
}

func (group *AppGroup) getResourceName() string {
	return group.Name
}
