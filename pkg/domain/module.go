package domain

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	"github.com/core-tools/hsu-echo/pkg/contract"
)

const id = "echo"

type echoModule struct {
	service1 contract.Contract1
	service2 contract.Contract2
}

func NewEchoModule(logger logging.Logger) moduletypes.Module {
	return &echoModule{
		service1: NewService1(logger),
		service2: NewService2(logger),
	}
}

func (m *echoModule) ID() moduletypes.ModuleID {
	return id
}

func (m *echoModule) SetServiceGatewayFactory(factory moduletypes.ServiceGatewayFactory) {
}

func (m *echoModule) ServiceHandlersMap() moduletypes.ServiceHandlersMap {
	return moduletypes.ServiceHandlersMap{
		"service1": m.service1,
		"service2": m.service2,
	}
}

func (m *echoModule) Start(ctx context.Context) error {
	return nil
}

func (m *echoModule) Stop(ctx context.Context) error {
	return nil
}
