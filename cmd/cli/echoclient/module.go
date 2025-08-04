package echoclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modules"
	"github.com/core-tools/hsu-echo/pkg/contract"
)

const id = "echoclient"

type echoclient struct {
	logger   logging.Logger
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewEchoClientModule(logger logging.Logger) modules.Module {
	return &echoclient{
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

func (m *echoclient) ID() string {
	return id
}

func (m *echoclient) Initialize(directClosureProvider modules.DirectClosureProvider) error {
	return nil
}

func (m *echoclient) Start(ctx context.Context, gatewayFactory modules.GatewayFactory) error {
	m.wg.Add(1)
	loop := &loop{
		gatewayFactory: gatewayFactory,
		logger:         m.logger,
		stopChan:       m.stopChan,
		wg:             &m.wg,
	}
	go loop.run()
	return nil
}

func (m *echoclient) Stop(ctx context.Context) error {
	close(m.stopChan)
	m.wg.Wait()
	return nil
}

type loop struct {
	gatewayFactory modules.GatewayFactory
	logger         logging.Logger
	stopChan       chan struct{}
	wg             *sync.WaitGroup
}

func (l *loop) run() {
	defer l.wg.Done()

	l.logger.Infof("Starting echo client loop")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.echo()
		case <-l.stopChan:
			l.logger.Infof("Echo client loop stopped")
			return
		}
	}
}

func (l *loop) echo() error {
	ctx := context.Background()

	echoGatewayRef, err := l.gatewayFactory.NewGateway(ctx, "echo", "")
	if err != nil {
		l.logger.Errorf("Failed to create gateway: %v", err)
		return err
	}

	echoGateway, ok := echoGatewayRef.(contract.Contract)
	if !ok {
		l.logger.Errorf("Failed to cast gateway to function")
		return fmt.Errorf("failed to cast gateway to function")
	}

	message := "Hello, World!"

	l.logger.Infof("Echo request: %s", message)

	echoResponse, err := echoGateway.Echo(ctx, message)
	if err != nil {
		l.logger.Errorf("Failed to echo: %v", err)
	} else {
		l.logger.Infof("Echo response: %s", echoResponse)
	}

	return nil
}
