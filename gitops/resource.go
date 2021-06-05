package gitops

import (
	"errors"
	"gosh/log"
	"os"
)

type ResourceList func() ([]*Resource, error)

type Resource interface {
	Read() error
	Exists() bool
	GetFilePath() string
	mapFromKapitanFile(f *kapitanFile)
	mapToKapitanFile() *kapitanFile
	isValid() bool
	getResourceType() string
	getResourceName() string
}

var (
	ResourceDoesNotExistErr  = errors.New("resource does not exist")
	ResourceAlreadyExistsErr = errors.New("resource already exist")
)

func Exists(resource Resource) bool {
	if f, err := os.Stat(resource.GetFilePath()); err == nil && !f.IsDir() {
		return true
	}
	return false
}

func Read(resource Resource) error {
	log.Tracef("Read %s with input: %+v", resource.getResourceName(), resource)
	if resource == nil || !resource.isValid() {
		return log.Errf(ValidationErr, "Invalid struct, use constructor to create one")
	}
	if !resource.Exists() {
		return log.Errf(ResourceDoesNotExistErr, "The %s '%s' does not exist", resource.getResourceType(), resource.getResourceName())
	}
	if f, err := ReadKapitanFile(resource.GetFilePath()); err == nil {
		resource.mapFromKapitanFile(f)
		log.Tracef("Read %s, result: %+v", resource.getResourceType(), resource)
		log.Infof("Read %s '%s'", resource.getResourceType(), resource.getResourceName())
		return nil
	} else {
		return log.Errf(err, "Could not read %s '%s' file", resource.getResourceType(), resource.getResourceName())
	}
}

func Update(resource Resource) error {
	log.Tracef("Update %s with input: %+v", resource.getResourceName(), resource)
	if resource == nil || !resource.isValid() {
		return log.Errf(ValidationErr, "Invalid struct, use constructor to create one")
	}
	if !resource.Exists() {
		return log.Errf(ResourceDoesNotExistErr, "The %s '%s' does not exist", resource.getResourceType(), resource.getResourceName())
	}
	f := resource.mapToKapitanFile()
	if err := WriteKapitanFile(resource.GetFilePath(), f); err == nil {
		log.Tracef("Updated %s, result: %+v", resource.getResourceType(), resource)
		log.Infof("Updated %s '%s'", resource.getResourceType(), resource.getResourceName())
		return nil
	} else {
		return log.Errf(err, "Could not update %s '%s'", resource.getResourceType(), resource.getResourceName())
	}
}
