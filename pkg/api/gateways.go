package echoapi

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduleapi"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/modulewiring"
	echocontract "github.com/core-tools/hsu-echo/pkg/api/contract"
	echogrpcapi "github.com/core-tools/hsu-echo/pkg/api/grpc"
)

func NewEchoServiceGateways(serviceConnector moduleapi.ServiceConnector, logger logging.Logger) echocontract.EchoServiceGateways {
	return &echoServiceGateways{
		targetModuleID:   moduletypes.ModuleID("echo"),
		serviceConnector: serviceConnector,
		logger:           logger,
	}
}

type echoServiceGateways struct {
	targetModuleID   moduletypes.ModuleID
	serviceConnector moduleapi.ServiceConnector
	serviceHandlers  echocontract.EchoServiceHandlers
	logger           logging.Logger
}

func (g *echoServiceGateways) ModuleID() moduletypes.ModuleID {
	return g.targetModuleID
}

func (g *echoServiceGateways) ServiceIDs() []moduletypes.ServiceID {
	return []moduletypes.ServiceID{
		moduletypes.ServiceID("service1"),
		moduletypes.ServiceID("service2"),
	}
}

func (g *echoServiceGateways) EnableDirectClosure(serviceHandlers echocontract.EchoServiceHandlers) {
	g.serviceHandlers = serviceHandlers
}

func (g *echoServiceGateways) GetService1(ctx context.Context, protocol moduletypes.Protocol) (echocontract.Service1, error) {
	serviceGatewayFactory := modulewiring.ServiceGatewayFactory[echocontract.Service1]{
		ModuleID:         g.targetModuleID,
		ServiceID:        moduletypes.ServiceID("service1"),
		ServiceConnector: g.serviceConnector,
		GatewayFactoryFuncs: modulewiring.GatewayFactoryFuncs[echocontract.Service1]{
			Direct: modulewiring.NewGatewayFactoryFuncDirect(g.serviceHandlers.Service1),
			GRPC:   echogrpcapi.NewGRPCGateway1,
			HTTP:   nil, // TODO: implement HTTP gateway
		},
		Logger: g.logger,
	}
	return serviceGatewayFactory.NewServiceGateway(ctx, protocol)
}

func (g *echoServiceGateways) GetService2(ctx context.Context, protocol moduletypes.Protocol) (echocontract.Service2, error) {
	serviceGatewayFactory := modulewiring.ServiceGatewayFactory[echocontract.Service2]{
		ModuleID:         g.targetModuleID,
		ServiceID:        moduletypes.ServiceID("service2"),
		ServiceConnector: g.serviceConnector,
		GatewayFactoryFuncs: modulewiring.GatewayFactoryFuncs[echocontract.Service2]{
			Direct: modulewiring.NewGatewayFactoryFuncDirect(g.serviceHandlers.Service2),
			GRPC:   nil, // TODO: implement echogrpcapi.NewGRPCGateway2,
			HTTP:   nil, // TODO: implement HTTP gateway
		},
		Logger: g.logger,
	}
	return serviceGatewayFactory.NewServiceGateway(ctx, protocol)
}

func EchoDirectClosureEnable(options modulewiring.DirectClosureEnableOptions[echocontract.EchoServiceGateways, echocontract.EchoServiceHandlers]) {
	serviceConnector := options.ServiceConnector
	serviceGateways := options.ServiceGateways
	serviceHandlers := options.ServiceHandlers
	serviceConnector.EnableDirectClosure(serviceGateways.ModuleID(), serviceGateways.ServiceIDs())
	serviceGateways.EnableDirectClosure(serviceHandlers)
}
