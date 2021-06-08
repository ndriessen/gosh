package gitops

import (
	"gosh/log"
)

//AppList An app list is a generic interface for Resource implementation that can give version information about a list of apps like Stage or Release
type AppList interface {
	//getResourceName Returns the name of the Resource this AppList is linked to
	getResourceName() string
	//GetVersions Return a map containing App names as keys and version numbers as values.
	//
	// If you specify either group or app, the list will be filtered by either the AppGroup or the App name.
	// If you specify both, filtering on the App name will be applied.
	GetVersions(group string, app string) map[string]string
	//GetArtifacts Return a map containing App names as keys and artifact URLs as values.
	//
	// If you specify either group or app, the list will be filtered by either the AppGroup or the App name.
	// If you specify both, filtering on the App name will be applied.
	//
	// the artifactType determines which artifact to return, see the App documentation about artifacts for more information
	GetArtifacts(group string, app string, artifactType string) (map[string]string, error)
	//versions this should return the complete map of App.GetArtifact names as keys and versions as values that this Resource contains
	//
	//Used to provide generic GetVersions and GetArtifacts implementations
	versions() map[string]string
	//UpdateVersion Updates the version for the app on the current Resource
	UpdateVersion(app string, version string) error
}

func GetVersions(list AppList, group string, app string) map[string]string {
	return filterList(list.versions(), group, app)
}

func GetArtifacts(list AppList, group string, app string, artifactType string) (artifacts map[string]string, err error) {
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

func getAppsToFilter(group string, app string) []string {
	if group != "" && app != "" {
		log.Warn("Both group and app filters are supplied... ignoring group filter")
	}
	var filtered []string
	if app != "" {
		filtered = append(filtered, app)
	} else {
		if group != "" {
			g := NewAppGroup(group)
			if err := g.Read(); err == nil && g.Exists() {
				for _, a := range g.Apps {
					filtered = append(filtered, a.Name)
				}
			}
		}
	}
	return filtered
}

func filterList(list map[string]string, group string, app string) map[string]string {
	filterAppNames := getAppsToFilter(group, app)
	if filterAppNames != nil && len(filterAppNames) > 0 {
		filtered := map[string]string{}
		for _, a := range filterAppNames {
			if v, exists := list[a]; exists {
				filtered[a] = v
			}
		}
		return filtered
	}
	return list
}
