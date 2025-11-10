package main

import (
	"fmt"
	"os"

	"github.com/core-tools/hsu-core/pkg/logging"
	sprintflogging "github.com/core-tools/hsu-core/pkg/logging/sprintf"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduleproto"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/modulewiring"

	_ "github.com/core-tools/hsu-echo/pkg/echoserver/echoserverwiring"

	flags "github.com/jessevdk/go-flags"
)

type flagOptions struct {
	Port int `long:"port" description:"port to listen on"`
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
				ID: "echo",
				Servers: []moduleproto.ServerID{
					"server-grpc",
				},
				Enabled: true,
			},
		},
		Runtime: modulewiring.RuntimeConfig{
			Servers: []modulewiring.ServerConfig{
				{
					ID:       "server-grpc",
					Protocol: moduletypes.ProtocolGRPC,
					Enabled:  true,
				},
			},
		},
	}

	err = modulewiring.RunWithConfig(config, logger)
	if err != nil {
		fmt.Printf("Failed to run module wiring: %v\n", err)
		os.Exit(1)
	}
}
