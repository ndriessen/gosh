package gosh_import

import (
	"encoding/json"
	"errors"
	"fmt"
	"gosh/gitops"
	"gosh/log"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TmImportPlugin struct {
}

func (p *TmImportPlugin) Name() string {
	return "trendminer"
}

func (p *TmImportPlugin) Import(apps bool, stages bool, releases bool, appTemplateName string) (err error) {
	if apps {
		err = importApps(appTemplateName)
	}
	if err == nil && stages {
		err = importVersions()
	}
	if err == nil && releases {
		err = importReleases()
	}
	return
}

var TrendMinerPlugin = &TmImportPlugin{}

func importVersions() error {
	log.Info("Starting import of versions")
	stages := map[string]string{"TESTED": "alpha", "PUBLISHED": "stable", "RELEASED": "released"}
	for oldStage, stageName := range stages {
		stage := gitops.NewStage(stageName)
		if !stage.Exists() {
			if err := stage.Create(); err != nil {
				log.Fatal(err, "Could not create stage %s", stageName)
			}
		}
		if err := stage.Read(); err != nil {
			log.Fatal(err, "Could not load stage %s", stageName)
		}
		log.Infof("Importing stage %s", stageName)
		if data, err := readUrl("http://versions.trendminer.net/versions/" + oldStage); err == nil {
			log.Tracef("Loaded data %s", data)
			versions := strings.Split(string(data), "\n")
			for _, line := range versions {
				parts := strings.Split(line, "=")
				if len(parts) == 2 {
					appName := strings.TrimSuffix(parts[0], ".version")
					if err = stage.UpdateVersion(appName, parts[1]); err == nil {
						log.Infof("Updated version: %s = %s", appName, parts[1])
					} else {
						_ = log.Errf(err, "Could not update version for %s = %s", appName, parts[1])
					}
				}
			}
		} else {
			return log.Errf(err, "Could not load versions from versions dashboard")
		}
	}
	return nil
}

func importReleases() error {
	log.Info("Starting import of releases, loading releases (this can take a while)")
	if data, err := readUrl("http://versions.trendminer.net/releases"); err == nil {
		log.Tracef("received response: %s", string(data))
		var p = new([]interface{})
		releaseData := make(map[string]map[string]interface{}, 0)
		if err = json.Unmarshal(data, p); err == nil {
			releases := make([]string, 0)
			//2021.R2(-11)
			years := []string{strconv.Itoa(time.Now().Year()), strconv.Itoa(time.Now().Year() - 1)}
			for _, release := range *p {
				name := release.(map[string]interface{})["name"].(string)
				releaseData[name] = release.(map[string]interface{})
				releases = append(releases, name)
				log.Tracef("Found release %s", name)
			}
			recentReleases := make([]string, 0)
			for _, year := range years {
				filtered := filterSlice(releases, year)
				for _, v := range filtered {
					sort.Strings(v)
					recentReleases = append(recentReleases, v[len(v)-2:]...)
				}
			}
			sort.Strings(recentReleases)
			for _, releaseName := range recentReleases {
				log.Infof("Importing release %s", releaseName)
				release := gitops.NewRelease(releaseName, gitops.ProductRelease)
				if release.Exists() && release.Read() != nil {
					log.Fatal(err, "Error importing release %s", releaseName)
				}
				var versions = releaseData[releaseName]["versions"]
				for _, v := range versions.([]interface{}) {
					version := v.(map[string]interface{})["version"].(string)
					app := v.(map[string]interface{})["project"].(map[string]interface{})["name"].(string)
					log.Tracef("%s = %s\n", app, version)
					release.Versions[app] = version
				}
				if release.Exists() {
					if err = release.Update(); err != nil {
						log.Fatal(err, "Error importing release %s", releaseName)
					}
				} else {
					if err = release.Create(); err != nil {
						log.Fatal(err, "Error importing release %s", releaseName)
					}
				}
				log.Infof("Successfully imported release %s", releaseName)
			}
			return nil
		} else {
			return log.Errf(err, "Error parsing JSON response")
		}
	} else {
		return log.Errf(err, "Could not load versions from versions dashboard")
	}
}

func filterSlice(slice []string, filter string) map[string][]string {
	re := regexp.MustCompile(`(\d{4}\.R\d+)(-\d+)?`)
	result := make(map[string][]string, 0)
	for _, e := range slice {
		if strings.Contains(strings.ToLower(e), strings.ToLower(filter)) {
			matches := re.FindStringSubmatch(e)
			if len(matches) > 1 {
				result[matches[1]] = append(result[matches[1]], e)
			}
		}
	}
	return result
}

func readUrl(url string) (data []byte, err error) {
	var client http.Client
	if resp, err := client.Get(url); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			if data, err = ioutil.ReadAll(resp.Body); err == nil {
				return data, nil
			} else {
				return nil, err
			}
		} else {
			return nil, errors.New("received " + fmt.Sprint(resp.StatusCode) + " response")
		}
	} else {
		return nil, err
	}
}

func importApps(appTemplateName string) error {
	if data, err := readUrl("http://versions.trendminer.net/projects"); err == nil {
		log.Tracef("received response: %s", string(data))
		var p = new([]interface{})
		if err = json.Unmarshal(data, p); err == nil {
			for _, project := range *p {
				projectType := strings.ToLower(project.(map[string]interface{})["type"].(string))
				name := project.(map[string]interface{})["name"].(string)
				group := gitops.NewAppGroup(projectType)
				app := gitops.NewApp(name, group)
				if app.Exists() {
					log.Warnf("App %s already exists... skipping", name)
					continue
				}
				app.Properties["groupId"] = strings.ReplaceAll(project.(map[string]interface{})["groupId"].(string), ".", "/")
				artifactId := project.(map[string]interface{})["artifactId"].(string)
				if projectType == "platform" && !strings.HasSuffix(artifactId, "dist") {
					artifactId += "-dist"
				}
				app.Properties["artifactId"] = artifactId
				if err = app.CreateFromTemplate(appTemplateName); err == nil {
					log.Infof("Imported app %s in group %s\n", name, projectType)
				} else {
					_ = log.Errf(err, "Error importing app %s", name)
				}
			}
			return err
		} else {
			return log.Errf(err, "Error importing apps")
		}
	} else {
		return log.Errf(err, "Error importing apps")
	}
}
