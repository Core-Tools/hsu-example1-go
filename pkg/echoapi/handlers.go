package echoapi

import (
	"github.com/core-tools/hsu-core/pkg/errors"
	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduleproto"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/modulewiring"
	"github.com/core-tools/hsu-echo/pkg/echocontract"
	"github.com/core-tools/hsu-echo/pkg/echogrpcapi"
)

func NewEchoHandlersRegistrar(protocolServers []moduleproto.ProtocolServer, logger logging.Logger) (modulewiring.HandlersRegistrar[echocontract.EchoServiceHandlers], error) {
	return &echoHandlersRegistrar{
		protocolServers: protocolServers,
		logger:          logger,
	}, nil
}

type echoHandlersRegistrar struct {
	protocolServers []moduleproto.ProtocolServer
	logger          logging.Logger
}

func (h *echoHandlersRegistrar) RegisterHandlers(handlers echocontract.EchoServiceHandlers) (modulewiring.ProtocolToServicesMap, error) {
	errCollection := errors.NewErrorCollection()
	errCollection.Add(h.registerService1(handlers.Service1))
	errCollection.Add(h.registerService2(handlers.Service2))
	err := errCollection.ToError()
	if err != nil {
		return nil, err
	}
	return map[moduletypes.Protocol][]moduletypes.ServiceID{
		moduletypes.ProtocolGRPC: []moduletypes.ServiceID{
			moduletypes.ServiceID("service1"),
		},
	}, nil
}

func (h *echoHandlersRegistrar) registerService1(service1 echocontract.Service1) error {
	registrarFuncs := modulewiring.RegisterHandlerFuncs[echocontract.Service1]{
		GRPC: echogrpcapi.RegisterGRPCHandler1,
		HTTP: nil, // TODO: implement HTTP handler
	}
	visitor := modulewiring.NewProtocolServerHandlersVisitor(service1, registrarFuncs, h.logger)
	return moduleproto.ApplyToProtocolServers(h.protocolServers, visitor)
}

func (h *echoHandlersRegistrar) registerService2(service2 echocontract.Service2) error {
	registrarFuncs := modulewiring.RegisterHandlerFuncs[echocontract.Service2]{
		GRPC: nil, // TODO: implement GRPC handler
		HTTP: nil, // TODO: implement HTTP handler
	}
	visitor := modulewiring.NewProtocolServerHandlersVisitor(service2, registrarFuncs, h.logger)
	return moduleproto.ApplyToProtocolServers(h.protocolServers, visitor)
}
