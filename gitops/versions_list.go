package gitops

import "gosh/log"

type VersionsList interface {
	GetVersions(group string, app string) map[string]string
	versions() map[string]string
}

func GetVersions(list VersionsList, group string, app string) map[string]string {
	if group != "" && app != "" {
		log.Warn("Both group and app filters are supplied... ignoring group filter")
	}
	var filterApps []string
	if app != "" {
		filterApps = append(filterApps, app)
	} else {
		if group != "" {
			g := NewAppGroup(group)
			if err := g.Read(); err == nil && g.Exists() {
				for _, a := range g.Apps {
					filterApps = append(filterApps, a.Name)
				}
			}
		}
	}
	if filterApps != nil && len(filterApps) > 0 {
		filtered := map[string]string{}
		for _, a := range filterApps {
			if v, exists := list.versions()[a]; exists {
				filtered[a] = v
			}
		}
		return filtered
	}
	return list.versions()
}
