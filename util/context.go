package util

import (
	"gosh/log"
	"os"
)

type GoshContext struct {
	WorkingDir string
}

var Context = &GoshContext{}

func init() {
	Context.WorkingDir, _ = os.Getwd()
	if wd, exists := os.LookupEnv("GOSH_WORKING_DIR"); exists {
		log.Debugf("Setting working dir from ENV: %s", os.ExpandEnv(wd))
		Context.WorkingDir = os.ExpandEnv(wd)
	}
}
