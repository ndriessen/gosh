package list

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"gosh/log"
	"gosh/util"
	"sort"
	"strings"
)

const DefaultOutputFormat = "yaml"

var (
	outputFormats              = map[string]OutputFormat{}
	UnsupportedOutputFormatErr = errors.New("unsupported output format")
	OutputFormatRenderErr      = errors.New("error rendering output format")
)

func newYamlOutputFormat() OutputFormat {
	return &YamlListOutputFormat{}
}

func newPropertiesOutputFormat() OutputFormat {
	return &PropertiesListOutputFormat{}
}

func init() {
	outputFormats["yaml"] = newYamlOutputFormat()
	outputFormats["properties"] = newPropertiesOutputFormat()
}

func Render(format string, list map[string]string, keySuffix string) (string, error) {
	if format == "" {
		format = util.Config.Output.DefaultFormat
		if format == "" {
			format = DefaultOutputFormat
		}
	}
	if format, exists := outputFormats[format]; exists {
		return format.Render(list, keySuffix)
	}
	return "", UnsupportedOutputFormatErr
}

type OutputFormat interface {
	Render(list map[string]string, keySuffix string) (string, error)
}

type YamlListOutputFormat struct{}

func (f *YamlListOutputFormat) Render(list map[string]string, keySuffix string) (string, error) {
	suffixed := make(map[string]string, 0)
	if keySuffix != "" {
		for k, v := range list {
			suffixed[fmt.Sprintf("%s.%s", k, keySuffix)] = v
		}
	} else {
		suffixed = list
	}
	if data, err := yaml.Marshal(suffixed); err == nil {
		return string(data), nil
	}
	return "", log.Errf(OutputFormatRenderErr, "Could not render YAML output for list %+v", list)
}

type PropertiesListOutputFormat struct{}

func (f *PropertiesListOutputFormat) Render(list map[string]string, keySuffix string) (string, error) {
	keys := make([]string, 0)
	for k := range list {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	data := ""
	for _, k := range keys {
		if keySuffix != "" {
			data += fmt.Sprintf("%s.%s=%s\n", k, keySuffix, list[k])
		} else {
			data += fmt.Sprintf("%s=%s\n", k, list[k])
		}
	}
	return strings.TrimSuffix(data, "\n"), nil
}
