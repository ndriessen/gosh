package gitops

import (
	"gosh/log"
	"gosh/util"
	"path/filepath"
	"strings"
)

const (
	releasesPath = kapitanClassesPath + "releases"
)

type Release struct {
	Type     ReleaseType
	Name     string
	Versions map[string]string
}

func NewRelease(name string, releaseType ReleaseType) *Release {
	return &Release{Name: name, Type: releaseType, Versions: map[string]string{}}
}

func (release *Release) Read() error {
	return Read(release)
}

func (release *Release) mapToKapitanFile() *kapitanFile {
	log.Tracef("Mapping release %s to kapitan file: %+v", release.Name, release)
	f := &kapitanFile{}
	props := f.Parameters[release.Name].(map[string]interface{})
	for key, value := range release.Versions {
		props[key] = map[string]string{"version": value}
	}
	log.Tracef("Mapped release %s to kapitan file, result: %+v", release.Name, f)
	return f
}

func (release *Release) mapFromKapitanFile(f *kapitanFile) {
	log.Tracef("Mapping release %s from kapitan file %+v", release.Name, f)
	release.Versions = make(map[string]string, 0)
	if properties, exists := f.Parameters[release.Name]; exists {
		for key, value := range properties.(map[interface{}]interface{}) {
			if version, exists := value.(map[interface{}]interface{})["version"]; exists {
				release.Versions[key.(string)] = version.(string)
			}
		}
	}
	log.Tracef("Mapped release %s from kapitan file, result: %+v", release.Name, release)
}

func (release *Release) Exists() bool {
	return Exists(release)
}

func (release *Release) GetFilePath() string {
	return filepath.Join(util.Context.WorkingDir, releasesPath, release.Type.String(), release.Name+kapitanFileExt)
}

func (release *Release) isValid() bool {
	return strings.TrimSpace(release.Name) != ""
}

func (release *Release) getResourceType() string {
	return "release"
}

func (release *Release) getResourceName() string {
	return release.Name
}
