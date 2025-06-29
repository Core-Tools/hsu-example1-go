package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	sprintfLogging "github.com/core-tools/hsu-core/pkg/logging/sprintf"

	coreControl "github.com/core-tools/hsu-core/pkg/control"
	coreDomain "github.com/core-tools/hsu-core/pkg/domain"
	coreLogging "github.com/core-tools/hsu-core/pkg/logging"

	echoControl "github.com/core-tools/hsu-echo/pkg/control"
	echoLogging "github.com/core-tools/hsu-echo/pkg/logging"

	flags "github.com/jessevdk/go-flags"
)

type flagOptions struct {
	ServerPath string `long:"server" description:"path to the server executable"`
	AttachPort int    `long:"port" description:"port to attach to the server"`
}

func logPrefix(module string) string {
	return fmt.Sprintf("module: %s-client , ", module)
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

	logger := sprintfLogging.NewStdSprintfLogger()

	logger.Infof("opts: %+v", opts)

	if opts.ServerPath == "" && opts.AttachPort == 0 {
		fmt.Println("Server path or attach port is required")
		os.Exit(1)
	}

	logger.Infof("Starting...")

	coreLogger := coreLogging.NewLogger(
		logPrefix("hsu-core"), coreLogging.LogFuncs{
			Debugf: logger.Debugf,
			Infof:  logger.Infof,
			Warnf:  logger.Warnf,
			Errorf: logger.Errorf,
		})
	echoLogger := echoLogging.NewLogger(
		logPrefix("hsu-echo"), echoLogging.LogFuncs{
			Debugf: logger.Debugf,
			Infof:  logger.Infof,
			Warnf:  logger.Warnf,
			Errorf: logger.Errorf,
		})

	coreConnectionOptions := coreControl.ConnectionOptions{
		ServerPath: opts.ServerPath,
		AttachPort: opts.AttachPort,
	}
	coreConnection, err := coreControl.NewConnection(coreConnectionOptions, coreLogger)
	if err != nil {
		logger.Errorf("Failed to create core connection: %v", err)
		return
	}

	coreClientGateway := coreControl.NewGRPCClientGateway(coreConnection.GRPC(), coreLogger)
	echoClientGateway := echoControl.NewGRPCClientGateway(coreConnection.GRPC(), echoLogger)

	ctx := context.Background()

	retryPingOptions := coreDomain.RetryPingOptions{
		RetryAttempts: 10,
		RetryInterval: 1 * time.Second,
	}
	err = coreDomain.RetryPing(ctx, coreClientGateway, retryPingOptions, coreLogger)
	if err != nil {
		logger.Errorf("Failed to ping Echo server: %v", err)
		return
	}

	n := 10
	message := "Hello, world!"
	{
		logger.Infof("Echos...")

		wg := sync.WaitGroup{}
		wg.Add(n)

		startTime := time.Now()

		echoFunc := func(i int) {
			defer wg.Done()

			logger.Infof("Echo %d request: %s", i, message)

			echoResponse, err := echoClientGateway.Echo(ctx, message)
			if err != nil {
				logger.Errorf("Failed to echo %d: %v", i, err)
				return
			}

			logger.Infof("Echo %d response: %s", i, echoResponse)
		}

		for i := 0; i < n; i++ {
			go echoFunc(i)
		}

		wg.Wait()

		elapsedTime := time.Since(startTime)
		logger.Infof("Total time: %.2fs", elapsedTime.Seconds())

		logger.Infof("Echos done")
	}

	logger.Infof("Done")
}
