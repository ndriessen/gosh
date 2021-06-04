package gitops

import "gosh/log"

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
