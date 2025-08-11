package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	osruntime "runtime"
	"syscall"

	"github.com/core-tools/hsu-core/pkg/logging"
	sprintflogging "github.com/core-tools/hsu-core/pkg/logging/sprintf"
	"github.com/core-tools/hsu-core/pkg/modules"
	"github.com/core-tools/hsu-core/pkg/runtime"
	grpcapi "github.com/core-tools/hsu-echo/pkg/api/grpc"
	"github.com/core-tools/hsu-echo/pkg/domain"
	"google.golang.org/grpc"

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
		fmt.Printf("Command line flags parsing failed: %v", err)
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

	if opts.Port == 0 {
		fmt.Println("Port is required")
		os.Exit(1)
	}

	logger.Infof("Starting...")

	moduleManager := modules.NewManager(logger)
	if err != nil {
		fmt.Println("Failed to create module manager")
		os.Exit(1)
	}

	moduleManager.RegisterModule("echo", domain.NewEchoSimpleModule(logger))

	moduleManager.ProvideHandlerRegistrar("echo", "", modules.HandlerConfig{
		GRPC: &modules.GRPCHandlerRegistrar{
			RegistrarFunc: func(grpcServiceRegistrar grpc.ServiceRegistrar, handler interface{}, logger logging.Logger) {
				grpcapi.RegisterGRPCHandler(grpcServiceRegistrar, handler, logger)
			},
		},
	})

	componentCtx := context.Background()
	operationCtx := componentCtx

	err = moduleManager.Initialize()
	if err != nil {
		fmt.Println("Failed to initialize module manager")
		os.Exit(1)
	}

	gatewayFactory := runtime.NewGatewayFactory(moduleManager, nil, logger)

	err = moduleManager.Start(operationCtx, gatewayFactory)
	if err != nil {
		fmt.Println("Failed to start module manager")
		os.Exit(1)
	}

	serverConfig := runtime.ServerConfig{
		GRPC: runtime.GRPCServerConfig{
			Port: opts.Port,
		},
	}
	server := runtime.NewServer(serverConfig, moduleManager, logger)

	server.Start(operationCtx)

	// Enable signal handling
	sig := make(chan os.Signal, 1)
	if osruntime.GOOS == "windows" {
		signal.Notify(sig) // Unix signals not implemented on Windows
	} else {
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	}

	logger.Infof("All components are ready, starting managed processes...")

	// Wait for graceful shutdown or timeout
	select {
	case receivedSignal := <-sig:
		logger.Infof("Master runner received signal: %v", receivedSignal)
		if osruntime.GOOS == "windows" {
			if receivedSignal != os.Interrupt {
				logger.Errorf("Wrong signal received: got %q, want %q\n", receivedSignal, os.Interrupt)
				os.Exit(42)
			}
		}
	case <-operationCtx.Done():
		logger.Infof("Master runner timed out")
	}

	logger.Infof("Ready to stop components...")

	// Stop runtime
	ctx := context.Background()
	server.Shutdown(ctx)
	moduleManager.Stop(ctx)

	logger.Infof("Done")
}
