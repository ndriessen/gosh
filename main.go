package main

import (
	"gosh/cmd"
	"gosh/util"
	"os"
)

var Version = "0.0.0"

func main() {
	util.Context.Version = Version
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
