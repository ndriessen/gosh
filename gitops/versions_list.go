package gitops

import (
	"gosh/log"
)

type VersionsList interface {
	getResourceName() string
	GetVersions(group string, app string) map[string]string
	GetArtifacts(group string, app string, artifactType string) (map[string]string, error)
	versions() map[string]string
}

func GetVersions(list VersionsList, group string, app string) map[string]string {
	return filterList(list.versions(), group, app)
}

func GetArtifacts(list VersionsList, group string, app string, artifactType string) (artifacts map[string]string, err error) {
	versions := filterList(list.versions(), group, app)
	artifacts = map[string]string{}
	for k, v := range versions {
		if app, err := FindApp(k); err == nil {
			if err = app.Read(); err == nil {
				if artifact, err := app.GetArtifact(list, v, artifactType); err == nil {
					artifacts[k] = artifact
				} else {
					return nil, log.Errf(err, "could not get artifact for app %s", app.Name)
				}
			}
		}
	}
	return
}
