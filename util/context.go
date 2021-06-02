package util

import "os"

type GoshContext struct {
	WorkingDir string
}

var Context = &GoshContext{}

func init() {
	if wd, exists := os.LookupEnv("GOSH_WORKING_DIR"); exists {
		Context.WorkingDir = os.ExpandEnv(wd)
	}
}
