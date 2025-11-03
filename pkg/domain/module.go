package echodomain

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	echocontract "github.com/core-tools/hsu-echo/pkg/api/contract"
)

type EchoServiceProvider interface {
}

func NewEchoModule(serviceProvider EchoServiceProvider, logger logging.Logger) (moduletypes.Module, echocontract.EchoServiceHandlers) {
	module := &echoModule{
		service1: NewHandler1(logger),
		service2: NewHandler2(logger),
	}
	return module,
		echocontract.EchoServiceHandlers{
			Service1: module.service1,
			Service2: module.service2,
		}
}

type echoModule struct {
	service1 echocontract.Service1
	service2 echocontract.Service2
}

func (m *echoModule) Start(ctx context.Context) error {
	return nil
}

func (m *echoModule) Stop(ctx context.Context) error {
	return nil
}
