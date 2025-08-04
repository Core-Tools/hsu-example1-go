package domain

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modules"
	"github.com/core-tools/hsu-echo/pkg/contract"
)

const id = "echo"

type echoSimple struct {
	handler contract.Contract
}

func NewEchoSimpleModule(logger logging.Logger) modules.Module {
	return &echoSimple{
		handler: NewSimpleHandler(logger),
	}
}

func (m *echoSimple) ID() string {
	return id
}

func (m *echoSimple) Initialize(directClosureProvider modules.DirectClosureProvider) error {
	directClosureProvider.ProvideDirectClosure(id, "", m.handler)

	return nil
}

func (m *echoSimple) Start(ctx context.Context, gatewayFactory modules.GatewayFactory) error {
	return nil
}

func (m *echoSimple) Stop(ctx context.Context) error {
	return nil
}
