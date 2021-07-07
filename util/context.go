package util

import (
	"errors"
	"gosh/log"
	"os"
)

type GoshContext struct {
	WorkingDir string
	Version    string
}

var Context = &GoshContext{}

func init() {
	Context.WorkingDir, _ = os.Getwd()
	if wd, exists := os.LookupEnv("GOSH_WORKING_DIR"); exists {
		Context.WorkingDir = os.ExpandEnv(wd)
		log.Debugf("Setting working dir from ENV: %s", Context.WorkingDir)
	}
	if wd, exists := os.LookupEnv("GOSH_WORK_DIR"); exists {
		Context.WorkingDir = os.ExpandEnv(wd)
		log.Debugf("Setting working dir from ENV: %s", Context.WorkingDir)
	}
	if _, err := os.Stat(Context.WorkingDir); err != nil {
		if _, ok := err.(*os.PathError); ok {
			if err = os.MkdirAll(Context.WorkingDir, 0755); err != nil {
				log.Fatal(errors.New("working Directory does not exist"), "Working directory does not exist")
			}
		}
	}
}
