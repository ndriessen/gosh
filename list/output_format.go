package list

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"gosh/log"
	"sort"
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

func Render(format string, list map[string]string) (string, error) {
	if format, exists := outputFormats[format]; exists {
		return format.Render(list)
	}
	return "", UnsupportedOutputFormatErr
}

type OutputFormat interface {
	Render(list map[string]string) (string, error)
}

type YamlListOutputFormat struct{}

func (f *YamlListOutputFormat) Render(list map[string]string) (string, error) {
	if data, err := yaml.Marshal(list); err == nil {
		return string(data), nil
	}
	return "", log.Errf(OutputFormatRenderErr, "Could not render YAML output for list %+v", list)
}

type PropertiesListOutputFormat struct{}

func (f *PropertiesListOutputFormat) Render(list map[string]string) (string, error) {
	keys := make([]string, 0)
	for k, _ := range list {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	data := ""
	for _, k := range keys {
		data += fmt.Sprintf("%s.version=%s\n", k, list[k])
	}
	return data, nil
}
