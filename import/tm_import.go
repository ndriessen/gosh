package gosh_import

import (
	"encoding/json"
	"errors"
	"fmt"
	"gosh/gitops"
	"gosh/log"
	"io/ioutil"
	"net/http"
	"strings"
)

type TmImportPlugin struct {
}

func (p *TmImportPlugin) Name() string {
	return "trendminer"
}

func (p *TmImportPlugin) Import() (err error) {
	err = importProjects()
	return
}

var TrendMinerPlugin = &TmImportPlugin{}

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
			return nil, errors.New("received " + string(resp.StatusCode) + " response")
		}
	} else {
		return nil, err
	}
}

func importProjects() error {
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
					fmt.Printf("App %s already exists... skipping\n", name)
					continue
				}
				app.Properties["groupId"] = project.(map[string]interface{})["groupId"].(string)
				app.Properties["artifactId"] = project.(map[string]interface{})["artifactId"].(string)
				if err = app.CreateFromTemplate("tm-app"); err == nil {
					log.Warnf("Imported app %s in group %s\n", name, projectType)
					fmt.Printf("Imported app %s in group %s\n", name, projectType)
				} else {
					_ = log.Errf(err, "Error importing app %s", name)
					fmt.Printf("Error importing app %s (%v)", name, err)
				}
			}
			return nil
		} else {
			return log.Errf(err, "Error importing apps")
		}
	} else {
		return log.Errf(err, "Error importing apps")
	}
}
