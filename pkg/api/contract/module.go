package echocontract

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
)

type EchoServiceHandlers struct {
	Service1 Service1
	Service2 Service2
}

type EchoServiceGateways interface {
	ModuleID() moduletypes.ModuleID
	ServiceIDs() []moduletypes.ServiceID
	EnableDirectClosure(serviceHandlers EchoServiceHandlers)
	GetService1(ctx context.Context, protocol moduletypes.Protocol) (Service1, error)
	GetService2(ctx context.Context, protocol moduletypes.Protocol) (Service2, error)
}
