package echoclient

import (
	"context"
	"sync"
	"time"

	"github.com/core-tools/hsu-core/pkg/errors"
	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	"github.com/core-tools/hsu-echo/pkg/contract"
)

const id = "echoclient"

type echoclient struct {
	logger   logging.Logger
	factory  moduletypes.ServiceGatewayFactory
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewEchoClientModule(logger logging.Logger) moduletypes.Module {
	return &echoclient{
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

func (m *echoclient) ID() moduletypes.ModuleID {
	return id
}

func (m *echoclient) SetServiceGatewayFactory(factory moduletypes.ServiceGatewayFactory) {
	m.factory = factory
}

func (m *echoclient) ServiceHandlersMap() moduletypes.ServiceHandlersMap {
	return nil
}

func (m *echoclient) Start(ctx context.Context) error {
	m.wg.Add(1)
	loop := &loop{
		logger:   m.logger,
		factory:  m.factory,
		stopChan: m.stopChan,
		wg:       &m.wg,
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
	factory  moduletypes.ServiceGatewayFactory
	logger   logging.Logger
	stopChan chan struct{}
	wg       *sync.WaitGroup
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

	targetModuleID := moduletypes.ModuleID("echo")
	targetServiceID := moduletypes.ServiceID("service1")
	targetProtocol := moduletypes.Protocol("auto")
	abstractService1, err := l.factory.NewServiceGateway(ctx,
		targetModuleID, targetServiceID, targetProtocol)
	if err != nil {
		l.logger.Errorf("Failed to create service gateway: %v", err)
		return err
	}

	typifiedService1, ok := abstractService1.(contract.Contract1)
	if !ok {
		l.logger.Errorf("Failed to cast service gateway: %T", abstractService1)
		return errors.NewDomainError(errors.ErrorTypeValidation, "failed to cast service gateway", nil).
			WithContext("service_gateway", abstractService1)
	}

	message := "Hello, World!"

	l.logger.Infof("Module client: Service1.Echo1 request: %s", message)

	echoResponse, err := typifiedService1.Echo1(ctx, message)
	if err != nil {
		l.logger.Errorf("Module client: Service1.Echo1 request failed: %v", err)
	} else {
		l.logger.Infof("Module client: Service1.Echo1 response: %s", echoResponse)
	}

	return nil
}
