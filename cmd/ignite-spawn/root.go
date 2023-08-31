package main

import (
	"fmt"
	"os"

	"github.com/save-abandoned-projects/ignite/pkg/logs"
	logflag "github.com/save-abandoned-projects/ignite/pkg/logs/flag"
	"github.com/save-abandoned-projects/ignite/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var logLevel = logrus.InfoLevel

// RunIgniteSpawn runs the root command for ignite-spawn
func RunIgniteSpawn() {
	fs := &pflag.FlagSet{
		Usage: usage,
	}

	addGlobalFlags(fs)
	util.GenericCheckErr(fs.Parse(os.Args[1:]))
	logs.Logger.SetLevel(logLevel)

	if len(fs.Args()) != 1 {
		usage()
	}

	util.GenericCheckErr(func() error {
		vm, err := decodeVM(fs.Args()[0])
		if err != nil {
			return err
		}

		return StartVM(vm)
	}())
}

func usage() {
	util.GenericCheckErr(fmt.Errorf("usage: ignite-spawn [--log-level <level>] <vm>"))
}

func addGlobalFlags(fs *pflag.FlagSet) {
	// TODO: Add a version flag
	logflag.LogLevelFlagVar(fs, &logLevel)
}
