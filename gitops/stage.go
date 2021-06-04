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
}

func NewStage(name string) *Stage {
	return &Stage{Name: name, Versions: map[string]string{}}
}

func (stage *Stage) Read() error {
	return Read(stage)
}

func (stage *Stage) mapToKapitanFile() *kapitanFile {
	log.Tracef("Mapping stage %s to kapitan file: %+v", stage.Name, stage)
	f := &kapitanFile{}
	props := f.Parameters[stage.Name].(map[string]interface{})
	for key, value := range stage.Versions {
		props[key] = value
	}
	log.Tracef("Mapped stage %s to kapitan file, result: %+v", stage.Name, f)
	return f
}

func (stage *Stage) mapFromKapitanFile(f *kapitanFile) {
	log.Tracef("Mapping stage %s from kapitan file %+v", stage.Name, f)
	stage.Versions = make(map[string]string, 0)
	if properties, exists := f.Parameters[stage.Name]; exists {
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
	return Exists(stage)
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
