package gitops

import (
	"errors"
	"gopkg.in/yaml.v2"
	"gosh/log"
	"io/ioutil"
	"os"
	"strings"
)

type kapitanFile struct {
	Classes    []string
	Parameters map[interface{}]interface{}
}

const (
	kapitanClassesPath = "inventory/classes/"
	kapitanFileExt     = ".yml"
)

var (
	missingArgumentErr  = errors.New("missing required argument")
	resourceNotFoundErr = errors.New("resource not found")
)

func fileExists(filePath string) bool {
	f, err := os.Stat(filePath)
	return err == nil && !f.IsDir()
}

func ReadKapitanFile(filePath string) (*kapitanFile, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, missingArgumentErr
	}
	if !fileExists(filePath) {
		return nil, resourceNotFoundErr
	}
	log.Debug("Reading kapitan resource", filePath)
	if data, err := ioutil.ReadFile(filePath); err == nil {
		var file = &kapitanFile{}
		if err = yaml.Unmarshal(data, file); err == nil {
			log.Trace("Read data", file)
			return file, nil
		} else {
			return nil, log.CheckErr(err, "could not parse kapitan resource", filePath)
		}
	} else {
		return nil, log.CheckErr(err, "error reading kapitan resource ", filePath)
	}
}

func WriteKapitanFile(filePath string, data *kapitanFile) error {
	if data == nil {
		return missingArgumentErr
	}
	if bytes, err := yaml.Marshal(data); err == nil {
		err = ioutil.WriteFile(filePath, bytes, 0644)
		return log.CheckErr(err, "error writing kapitan resource", filePath)
	} else {
		return log.CheckErr(err, "error converting data to YAML", data)
	}
}
