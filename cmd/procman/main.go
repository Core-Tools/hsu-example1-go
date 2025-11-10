package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/core-tools/hsu-core/pkg/logging"
	sprintflogging "github.com/core-tools/hsu-core/pkg/logging/sprintf"
	"github.com/core-tools/hsu-core/pkg/managedprocess"
	"github.com/core-tools/hsu-core/pkg/managedprocess/processcontrol"
	"github.com/core-tools/hsu-core/pkg/process"
	"github.com/core-tools/hsu-core/pkg/processmanagement"

	flags "github.com/jessevdk/go-flags"
)

type flagOptions struct {
	EchoServerPath string `long:"echo-server-path" description:"path to the echo server executable"`
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

	if opts.EchoServerPath == "" {
		fmt.Printf("Echo server path is required\n")
		os.Exit(1)
	}

	logger.Infof("Starting...")

	processManagerOptions := processmanagement.ProcessManagerOptions{}
	processManager := processmanagement.NewProcessManager(processManagerOptions, logger)

	processOptionsArr := make([]managedprocess.ProcessOptions, 0)
	{
		unit := &managedprocess.IntegratedManagedProcessConfig{
			Metadata: managedprocess.ProcessMetadata{
				Name: "Echo Server",
			},
			Control: processcontrol.ManagedProcessControlConfig{
				Execution: process.ExecutionConfig{
					ExecutablePath: opts.EchoServerPath,
				},
			},
		}
		processOptions := managedprocess.NewIntegratedManagedProcessOptions("echo", unit, logger)
		processOptionsArr = append(processOptionsArr, processOptions)
	}

	for _, processOptions := range processOptionsArr {
		err = processManager.AddProcess(processOptions)
		if err != nil {
			fmt.Printf("Failed to add worker: %v\n", err)
			os.Exit(1)
		}
	}

	componentCtx := context.Background()
	operationCtx := componentCtx

	err = processManager.Start(operationCtx)
	if err != nil {
		fmt.Printf("Failed to start process manager: %v\n", err)
		os.Exit(1)
	}

	// Enable signal handling
	sig := make(chan os.Signal, 1)
	if runtime.GOOS == "windows" {
		signal.Notify(sig) // Unix signals not implemented on Windows
	} else {
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	}

	logger.Infof("All components are ready, starting managed processes...")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Start all managed processes (lifecycle phase)
		for _, processOptions := range processOptionsArr {
			err := processManager.StartProcess(componentCtx, processOptions.ID())
			if err != nil {
				logger.Errorf("Failed to start process %s: %v", processOptions.ID(), err)
				// Continue with other managed processes rather than failing completely
				continue
			}
			logger.Infof("Started process: %s", processOptions.ID())
		}

		logger.Infof("All managed processes started, process manager is fully operational")
	}()

	// Wait for graceful shutdown or timeout
	select {
	case receivedSignal := <-sig:
		logger.Infof("Process manager runner received signal: %v", receivedSignal)
		if runtime.GOOS == "windows" {
			if receivedSignal != os.Interrupt {
				logger.Errorf("Wrong signal received: got %q, want %q\n", receivedSignal, os.Interrupt)
				os.Exit(42)
			}
		}
	case <-operationCtx.Done():
		logger.Infof("Process manager runner timed out")
	}

	logger.Infof("Waiting for managed processes start to finish...")

	// Wait for starting managed processes to finish
	wg.Wait()

	logger.Infof("Ready to stop components...")

	// Stop process manager
	ctx := context.Background()
	err = processManager.Stop(ctx)
	if err != nil {
		fmt.Printf("Failed to stop process manager: %v\n", err)
		os.Exit(1)
	}

	logger.Infof("Done")
}
