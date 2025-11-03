package main

import (
	"context"
	"fmt"
	"os"

	"github.com/core-tools/hsu-core/pkg/logging"
	sprintflogging "github.com/core-tools/hsu-core/pkg/logging/sprintf"
	"github.com/core-tools/hsu-core/pkg/modulemanagement"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduleapi"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	"github.com/core-tools/hsu-echo/cmd/cli/echoclient"
	grpcapi "github.com/core-tools/hsu-echo/pkg/api/grpc"

	flags "github.com/jessevdk/go-flags"
)

type flagOptions struct {
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

	componentCtx := context.Background()
	operationCtx := componentCtx

	modules := []moduletypes.Module{
		echoclient.NewEchoClientModule(logger),
	}
	moduleGatewayConfigs := moduleapi.ModuleGatewaysConfigMap{
		"echo": []moduleapi.ServiceGatewayConfig{
			{
				ServiceID:          "service1",
				Protocol:           moduletypes.ProtocolGRPC,
				GatewayFactoryFunc: grpcapi.NewGRPCGateway1,
			},
		},
	}

	runtimeOptions := modulemanagement.RuntimeOptions{
		Modules:               modules,
		ModuleGatewaysConfigs: moduleGatewayConfigs,
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
