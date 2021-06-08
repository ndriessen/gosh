package gitops

import (
	"errors"
	"gosh/log"
	"gosh/util"
	"path/filepath"
	"strings"
)

const (
	releasesPath = kapitanClassesPath + "releases"
)

var (
	InvalidFullReleaseNameErr = errors.New("invalid release name, must be 'type/name' and type 'stage' is reserved")
)

type Release struct {
	Type     ReleaseType
	Name     string
	Versions map[string]string
	_read    bool
}

func (release *Release) initialized() bool {
	return release._read
}

func (release *Release) setInitialized() {
	release._read = true
}

func NewRelease(name string, releaseType ReleaseType) *Release {
	return &Release{Name: name, Type: releaseType, Versions: map[string]string{}}
}

func NewReleaseFromFullName(fullName string) (*Release, error) {
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return nil, InvalidFullReleaseNameErr
	}
	name := parts[1]
	if t, err := NewReleaseType(parts[0]); err == nil {
		if t == StageRelease {
			return nil, InvalidFullReleaseNameErr
		}
		return NewRelease(name, t), nil
	} else {
		return nil, log.Errf(err, "Unsupported release type", parts[0])
	}
}

func (release *Release) Create() error {
	return create(release)
}

func (release *Release) Read() error {
	return read(release)
}

func (release *Release) Update() error {
	return update(release)
}

func (release *Release) UpdateVersion(appName string, version string) error {
	if err := release.Read(); err == nil {
		if app, err := FindApp(appName); err == nil {
			release.Versions[app.Name] = version
			return release.Update()
		} else {
			return log.Errf(ResourceDoesNotExistErr, "the app with name %s does not exist", appName)
		}
	} else {
		return err
	}

}

func (release *Release) mapToKapitanFile() *kapitanFile {
	log.Tracef("Mapping release %s to kapitan file: %+v", release.Name, release)
	f := newKapitanFile()
	f.Parameters[release.Name] = make(map[string]interface{}, 0)
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

func (release *Release) versions() map[string]string {
	return release.Versions
}

func (release *Release) GetVersions(group string, app string) map[string]string {
	return GetVersions(release, group, app)
}

func (release *Release) GetArtifacts(group string, app string, artifactType string) (map[string]string, error) {
	return GetArtifacts(release, group, app, artifactType)
}

func (release *Release) Exists() bool {
	return exists(release)
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
