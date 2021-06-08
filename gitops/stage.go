package gitops

import (
	"gosh/log"
	"gosh/util"
	"path/filepath"
	"strings"
)

const (
	stagesPath = kapitanClassesPath + "stages"
)

type Stage struct {
	Name     string
	Versions map[string]string
	_read    bool
}

func (stage *Stage) initialized() bool {
	return stage._read
}

func (stage *Stage) setInitialized() {
	stage._read = true
}

func NewStage(name string) *Stage {
	return &Stage{Name: name, Versions: map[string]string{}}
}

func (stage *Stage) UpdateVersion(appName string, version string) error {
	if err := stage.Read(); err == nil {
		if app, err := FindApp(appName); err == nil {
			stage.Versions[app.Name] = version
			if err = stage.Update(); err == nil {
				release := NewRelease(stage.Name, StageRelease)
				if err = release.Read(); err == nil {
					return release.UpdateVersion(appName, version)
				} else {
					return log.Errf(err, "Could not update stage release %s", release.Name)
				}
			} else {
				return log.Errf(err, "Could not update stage %s", stage.Name)
			}
		} else {
			return log.Errf(ResourceDoesNotExistErr, "app with name %s does not exist", appName)
		}
	} else {
		return err
	}

}

func (stage *Stage) Create() error {
	if err := create(stage); err == nil {
		stageRelease := NewRelease(stage.Name, StageRelease)
		stageRelease.Versions = stage.Versions
		if !stageRelease.Exists() {
			if err = stageRelease.Create(); err != nil {
				return log.Errf(err, "Could not create associated stage release %s", stage.Name)
			}
		} else {
			//to sync versions
			_ = stageRelease.Update()
		}
		log.Infof("Created stage '%s'", stage.Name)
		return nil
	} else {
		return log.Errf(err, "Could not create stage %s", stage.Name)
	}
}

func (stage *Stage) Read() error {
	return read(stage)
}

func (stage *Stage) Update() error {
	return update(stage)
}

func (stage *Stage) mapToKapitanFile() *kapitanFile {
	log.Tracef("Mapping stage %s to kapitan file: %+v", stage.Name, stage)
	f := newKapitanFile()
	f.Parameters[stage.Name] = make(map[string]string, 0)
	props := f.Parameters[stage.Name].(map[string]string)
	for key, value := range stage.Versions {
		props[key] = value
	}
	log.Tracef("Mapped stage %s to kapitan file, result: %+v", stage.Name, f)
	return f
}

func (stage *Stage) mapFromKapitanFile(f *kapitanFile) {
	log.Tracef("Mapping stage %s from kapitan file %+v", stage.Name, f)
	stage.Versions = make(map[string]string, 0)
	if properties, exists := f.Parameters[stage.Name]; exists && properties != nil {
		for key, value := range properties.(map[interface{}]interface{}) {
			stage.Versions[key.(string)] = value.(string)
		}
	}
	log.Tracef("Mapped stage %s from kapitan file, result: %+v", stage.Name, stage)
}

func (stage *Stage) versions() map[string]string {
	return stage.Versions
}

func (stage *Stage) GetVersions(group string, app string) map[string]string {
	return GetVersions(stage, group, app)
}

func (stage *Stage) GetArtifacts(group string, app string, artifactType string) (map[string]string, error) {
	return GetArtifacts(stage, group, app, artifactType)
}

func (stage *Stage) Exists() bool {
	return exists(stage)
}

func (stage *Stage) GetFilePath() string {
	return filepath.Join(util.Context.WorkingDir, stagesPath, stage.Name+kapitanFileExt)
}

func (stage *Stage) isValid() bool {
	return strings.TrimSpace(stage.Name) != ""
}

func (stage *Stage) getResourceType() string {
	return "stage"
}

func (stage *Stage) getResourceName() string {
	return stage.Name
}
