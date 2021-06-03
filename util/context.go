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
	if wd, exists := os.LookupEnv("GOSH_WORKING_DIR"); exists {
		log.Infof("Setting working dir from ENV: %s", os.ExpandEnv(wd))
		Context.WorkingDir = os.ExpandEnv(wd)
	}
}
