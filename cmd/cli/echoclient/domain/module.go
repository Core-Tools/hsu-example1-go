package echoclientdomain

import (
	"context"
	"sync"
	"time"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	echocontract "github.com/core-tools/hsu-echo/pkg/api/contract"
)

type EchoClientServiceProvider interface {
	Echo() echocontract.EchoServiceGateways
}

func NewEchoClientModule(serviceProvider EchoClientServiceProvider, logger logging.Logger) (moduletypes.Module, moduletypes.EmptyServiceHandlers) {
	module := &echoclientModule{
		logger:          logger,
		serviceProvider: serviceProvider,
		stopChan:        make(chan struct{}),
	}
	return module, moduletypes.EmptyServiceHandlers{}
}

type echoclientModule struct {
	logger          logging.Logger
	serviceProvider EchoClientServiceProvider
	stopChan        chan struct{}
	wg              sync.WaitGroup
}

func (m *echoclientModule) Start(ctx context.Context) error {
	m.wg.Add(1)
	loop := &loop{
		logger:          m.logger,
		serviceProvider: m.serviceProvider,
		stopChan:        m.stopChan,
		wg:              &m.wg,
	}
	go loop.run()
	return nil
}

func (m *echoclientModule) Stop(ctx context.Context) error {
	close(m.stopChan)
	m.wg.Wait()
	return nil
}

type loop struct {
	serviceProvider EchoClientServiceProvider
	logger          logging.Logger
	stopChan        chan struct{}
	wg              *sync.WaitGroup
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

	echoService1, err := l.serviceProvider.Echo().GetService1(ctx, moduletypes.Protocol("auto"))
	if err != nil {
		l.logger.Errorf("Failed to get service1: %v", err)
		return err
	}

	message := "Hello, World!"

	l.logger.Infof("Module client: Service1.Echo1 request: %s", message)

	echoResponse, err := echoService1.Echo1(ctx, message)
	if err != nil {
		l.logger.Errorf("Module client: Service1.Echo1 request failed: %v", err)
	} else {
		l.logger.Infof("Module client: Service1.Echo1 response: %s", echoResponse)
	}

	return nil
}
