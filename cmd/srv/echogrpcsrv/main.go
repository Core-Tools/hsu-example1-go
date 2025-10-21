package main

import (
	"context"
	"fmt"
	"os"

	"github.com/core-tools/hsu-core/pkg/logging"
	sprintflogging "github.com/core-tools/hsu-core/pkg/logging/sprintf"
	"github.com/core-tools/hsu-core/pkg/modulemanagement"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduleapi"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduleproto"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	grpcapi "github.com/core-tools/hsu-echo/pkg/api/grpc"
	"github.com/core-tools/hsu-echo/pkg/domain"

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

	logger.Infof("Starting...")

	// 'echo' server module, remotely-accessible via gRPC

	componentCtx := context.Background()
	operationCtx := componentCtx

	module1 := domain.NewEchoModule(logger)
	modules := []moduletypes.Module{
		module1,
	}
	server1ID := moduleproto.ServerID("echogrpcsrv")
	serverOptions := moduleproto.ServerOptionsList{
		moduleproto.ServerOptions{
			ServerID: server1ID,
			Protocol: moduletypes.ProtocolGRPC,
			ProtocolOptions: moduleproto.GRPCServerOptions{
				Port: opts.Port, // this is only a config hint, could be 0 for dynamic port allocation
			},
		},
	}
	moduleHandlersConfigs := []moduleapi.ModuleHandlersConfig{
		{
			ModuleID:              module1.ID(),
			ServerID:              server1ID,
			Protocol:              moduletypes.ProtocolGRPC,
			HandlersRegistrarFunc: grpcapi.RegisterHandlers,
		},
	}

	runtimeOptions := modulemanagement.RuntimeOptions{
		Modules:               modules,
		ServerOptions:         serverOptions,
		ModuleHandlersConfigs: moduleHandlersConfigs,
		Logger:                logger,
	}

	runtime, err := modulemanagement.NewRuntime(runtimeOptions)
	if err != nil {
		fmt.Printf("Failed to create runtime: %v\n", err)
		os.Exit(1)
	}

	err = runtime.Start(operationCtx)
	if err != nil {
		fmt.Printf("Failed to start runtime: %v\n", err)
		os.Exit(1)
	}

	logger.Infof("Runtime is ready")

	modulemanagement.WaitSignals(operationCtx, logger)

	logger.Infof("About to stop runtime...")

	// Stop runtime
	ctx := context.Background()
	err = runtime.Stop(ctx)
	if err != nil {
		fmt.Printf("Failed to stop runtime: %v\n", err)
		os.Exit(1)
	}

	logger.Infof("Done")
}
