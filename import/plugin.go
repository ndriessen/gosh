package gosh_import

import (
	"errors"
	"gosh/log"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

const pluginSymbolName = "Plugin"

var (
	PluginPath        = os.ExpandEnv("$HOME/.gosh/plugins")
	PluginNotFoundErr = errors.New("plugin not found")
	BundledPlugins    = make(map[string]ImportPlugin, 0)
)

func init() {
	var t interface{} = TrendMinerPlugin
	if p, ok := t.(ImportPlugin); ok {
		BundledPlugins[p.Name()] = p
	} else {
		log.Warnf("Could not load bundled plugin 'trendminer', not implementing ImportPlugin{} correctly")
	}
}

type ImportPlugin interface {
	Name() string
	Import(apps bool, stages bool, releases bool) error
}

func Import(pluginName string, apps bool, stages bool, releases bool) (err error) {
	if p, exists := BundledPlugins[pluginName]; exists {
		return p.Import(apps, stages, releases)
	}
	if p, err := loadPlugin(pluginName); err == nil {
		err = runPlugin(pluginName, p, apps, stages, releases)
	}
	return
}

func runPlugin(pluginName string, p ImportPlugin, apps bool, stages bool, releases bool) error {
	log.Infof("Running import plugin %s", pluginName)
	if err := p.Import(apps, stages, releases); err == nil {
		log.Infof("Import with plugin %s successful", pluginName)
	} else {
		return log.Errf(err, "Import with plugin %s failed", pluginName)
	}
	return nil
}

func ListPlugins() ([]string, error) {
	list := make([]string, 0)
	if info, err := os.Stat(PluginPath); err == nil && info.IsDir() {
		if dir, err := os.ReadDir(PluginPath); err == nil {
			for _, entry := range dir {
				list = append(list, strings.TrimSuffix(entry.Name(), ".so"))
			}
		}
	}
	return list, nil
}

func loadPlugin(pluginName string) (ImportPlugin, error) {
	pluginPath := filepath.Join(PluginPath, pluginName+".so")
	if _, err := os.Stat(pluginPath); err != nil {
		return nil, log.Errf(PluginNotFoundErr, "plugin %s not found at path %s", pluginName, pluginPath)
	}
	if p, err := plugin.Open(pluginPath); err == nil {
		if symbol, err := p.Lookup(pluginSymbolName); err == nil {
			if importPlugin, ok := symbol.(ImportPlugin); ok {
				return importPlugin, nil
			} else {
				return nil, log.Errf(err, "plugin %s does not implement the gosh_import.ImportPlugin interface", pluginName)
			}
		} else {
			return nil, log.Errf(err, "plugin %s does not expose a symbol named 'TrendMinerPlugin'", pluginName)
		}
	} else {
		return nil, log.Errf(err, "error loading plugin: %s", pluginName)
	}
}
