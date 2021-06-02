package gitops

import (
	"errors"
	"gosh/log"
	"gosh/util"
	"os"
	"path/filepath"
	"strings"
)

const (
	appGroupPath = kapitanClassesPath + "app"
	appPrefix    = "app."
)

var (
	AppGroupAlreadyExistsErr = errors.New("app group already exists")
	AppGroupDoesNotExistErr  = errors.New("app group does not exist")
)

type AppGroup struct {
	Name string
	Apps []*App
}

// NewAppGroup /** Constructor for an AppGroup type
func NewAppGroup(name string, apps ...*App) *AppGroup {
	return &AppGroup{Name: name, Apps: apps}
}

func (group *AppGroup) Create() error {
	log.Trace("Create app group with input: %i", group)
	if !validStruct(group) {
		return log.Err(ValidationErr, "Invalid app group struct, use NewAppGroup() to create one")
	}
	if group.Exists() {
		return log.Errf(AppGroupAlreadyExistsErr, "The app group '%s' already exists", group.Name)
	}
	if err := os.MkdirAll(group.GetFolderPath(), 0755); err == nil {
		f := group.mapToKapitanFile()
		if err = WriteKapitanFile(group.GetFilePath(), f); err == nil {
			log.Infof("Created app group '%s", group.Name)
			return nil
		} else {
			return log.Errf(err, "Error writing app group file '%s'", group.GetFilePath())
		}
	} else {
		return log.Err(err, "Error creating app group folder '%s'", group.GetFolderPath())
	}
}

func (group *AppGroup) mapToKapitanFile() *kapitanFile {
	log.Tracef("Mapping group %s to kapitan file: %i", group.Name, group)
	f := &kapitanFile{}
	for _, app := range group.Apps {
		f.Classes = append(f.Classes, "app."+group.Name+"."+app.Name)
	}
	log.Tracef("Mapped group %s to kapitan file, result: %i", group.Name, f)
	return f
}

func (group *AppGroup) mapFromKapitanFile(f *kapitanFile) {
	log.Tracef("Mapping group %s from kapitan file %i", group.Name, f)
	for _, class := range f.Classes {
		app := &App{
			Name:  strings.TrimPrefix(class, appPrefix+group.Name+"."),
			group: group,
		}
		//if properties, exists := f.Parameters[app.Name]; exists {
		//	for key, value := range properties.(map[string]string) {
		//		app.Properties[key] = value
		//	}
		//}
		group.Apps = append(group.Apps, app)
	}
	log.Tracef("Mapped group %s from kapitan file, result: %i", group.Name, group)
}

func (group *AppGroup) Read() error {
	log.Tracef("Read app group with input: %i", group)
	if !validStruct(group) {
		return log.Err(ValidationErr, "Invalid app group struct, use NewAppGroup() to create one")
	}
	if !group.Exists() {
		return log.Errf(AppGroupDoesNotExistErr, "The app group '%s' does not exist", group.Name)
	}
	if f, err := ReadKapitanFile(group.GetFilePath()); err == nil {
		group.mapFromKapitanFile(f)
		log.Tracef("Read app group, result: %i", group)
		log.Infof("Read app group '%s'", group.Name)
		return nil
	} else {
		return log.Errf(err, "Could not read app group '%s' file", group.Name)
	}
}

func (group *AppGroup) Update() error {
	return nil
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

func validStruct(group *AppGroup) bool {
	return group != nil && strings.TrimSpace(group.Name) != ""
}
