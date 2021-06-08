package gitops

import (
	"errors"
	"gosh/log"
	"os"
)

type Resource interface {
	Create() error
	Read() error
	Update() error
	Exists() bool
	GetFilePath() string
	mapFromKapitanFile(f *kapitanFile)
	mapToKapitanFile() *kapitanFile
	isValid() bool
	getResourceType() string
	getResourceName() string
	initialized() bool
	setInitialized()
}

var (
	ResourceDoesNotExistErr          = errors.New("resource does not exist")
	ResourceAlreadyExistsErr         = errors.New("resource already exist")
	ResourceUpdatedWithoutReadingErr = errors.New("the resource was updated before being read from disk, this might lead to data loss")
)

func exists(resource Resource) bool {
	if f, err := os.Stat(resource.GetFilePath()); err == nil && !f.IsDir() {
		return true
	}
	return false
}

func create(resource Resource) error {
	log.Tracef("Create resource %s with input: %+v", resource.getResourceType(), resource)
	if resource == nil || !resource.isValid() {
		return log.Err(ValidationErr, "Invalid resource struct, use constructor to create one")
	}
	if resource.Exists() {
		return log.Errf(ResourceAlreadyExistsErr, "The %s '%s' already exists", resource.getResourceType(), resource.getResourceName())
	}
	f := resource.mapToKapitanFile()
	if err := WriteKapitanFile(resource.GetFilePath(), f); err == nil {
		log.Infof("Created %s '%s", resource.getResourceType(), resource.getResourceName())
		return nil
	} else {
		return log.Errf(err, "Error writing %s file '%s'", resource.getResourceType(), resource.GetFilePath())
	}
}

func read(resource Resource) error {
	if resource.initialized() {
		log.Tracef("Resource already read, skipping")
		return nil
	}
	log.Tracef("Read %s with input: %+v", resource.getResourceName(), resource)
	if resource == nil || !resource.isValid() {
		return log.Errf(ValidationErr, "Invalid struct, use constructor to create one")
	}
	if !resource.Exists() {
		return log.Errf(ResourceDoesNotExistErr, "The %s '%s' does not exist", resource.getResourceType(), resource.getResourceName())
	}
	if f, err := ReadKapitanFile(resource.GetFilePath()); err == nil {
		resource.mapFromKapitanFile(f)
		resource.setInitialized()
		log.Tracef("Read %s, result: %+v", resource.getResourceType(), resource)
		log.Infof("Read %s '%s'", resource.getResourceType(), resource.getResourceName())
		return nil
	} else {
		return log.Errf(err, "Could not read %s '%s' file", resource.getResourceType(), resource.getResourceName())
	}
}

func update(resource Resource) error {
	log.Tracef("Update %s with input: %+v", resource.getResourceName(), resource)
	if resource == nil || !resource.isValid() {
		return log.Errf(ValidationErr, "Invalid struct, use constructor to create one")
	}
	if !resource.Exists() {
		return log.Errf(ResourceDoesNotExistErr, "The %s '%s' does not exist", resource.getResourceType(), resource.getResourceName())
	}
	if !resource.initialized() {
		return log.Errf(ResourceUpdatedWithoutReadingErr, "Read the %s '%s' before updating it", resource.getResourceType(), resource.getResourceName())
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
