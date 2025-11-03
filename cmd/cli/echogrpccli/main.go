package main

import (
	"fmt"
	"os"

	"github.com/core-tools/hsu-core/pkg/logging"
	sprintflogging "github.com/core-tools/hsu-core/pkg/logging/sprintf"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/modulewiring"

	_ "github.com/core-tools/hsu-example1-go/cmd/cli/echoclient/app"

	flags "github.com/jessevdk/go-flags"
)

type flagOptions struct {
	//ConfigFile string `long:"config" description:"path to the config file"`
}

func main() {
	var opts flagOptions
	var argv []string = os.Args[1:]
	var parser = flags.NewParser(&opts, flags.HelpFlag)
	var err error
	_, err = parser.ParseArgs(argv)
	if err != nil {
		fmt.Printf("Command line flags parsing failed: %v\n", err)
		os.Exit(1)
	}

	sprintfLogger := sprintflogging.NewStdSprintfLogger()

	logger := logging.NewLogger(
		"",
		logging.LogFuncs{
			Debugf: sprintfLogger.Debugf,
			Infof:  sprintfLogger.Infof,
			Warnf:  sprintfLogger.Warnf,
			Errorf: sprintfLogger.Errorf,
		},
	)

	logger.Infof("opts: %+v", opts)

	config := &modulewiring.Config{
		Modules: []modulewiring.ModuleConfig{
			{
				ID:      "echo-client",
				Enabled: true,
			},
		},
	}

	err = modulewiring.RunWithConfig(config, logger)
	if err != nil {
		logger.Errorf("Failed to run module wiring: %v", err)
		os.Exit(1)
	}
}
