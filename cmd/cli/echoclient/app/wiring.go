package echoclientapp

import (
	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduleapi"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/modulewiring"
	echoapi "github.com/core-tools/hsu-echo/pkg/api"
	echocontract "github.com/core-tools/hsu-echo/pkg/api/contract"
	echoclientdomain "github.com/core-tools/hsu-example1-go/cmd/cli/echoclient/domain"
)

func init() {
	// Self-register in global registry
	moduleDesc := modulewiring.ModuleDescriptor[
		echoclientdomain.EchoClientServiceProvider,
		moduletypes.EmptyServiceGateways,
		moduletypes.EmptyServiceHandlers,
	]{
		ServiceProviderFactoryFunc:   NewEchoClientServiceProvider,
		ModuleFactoryFunc:            echoclientdomain.NewEchoClientModule,
		HandlersRegistrarFactoryFunc: nil,
		DirectClosureEnableFunc:      nil,
	}
	modulewiring.RegisterModule("echo-client", moduleDesc)
}

func NewEchoClientServiceProvider(serviceConnector moduleapi.ServiceConnector, logger logging.Logger) moduletypes.ServiceProviderHandle {
	echoServiceGateways := echoapi.NewEchoServiceGateways(serviceConnector, logger)
	serviceGatewaysMap := moduletypes.ServiceGatewaysMap{
		echoServiceGateways.ModuleID(): echoServiceGateways,
	}
	serviceProvider := &echoClientServiceProvider{
		echoServiceGateways: echoServiceGateways,
	}
	return moduletypes.ServiceProviderHandle{
		ServiceProvider:    serviceProvider,
		ServiceGatewaysMap: serviceGatewaysMap,
	}
}

type echoClientServiceProvider struct {
	echoServiceGateways echocontract.EchoServiceGateways
}

func (p *echoClientServiceProvider) Echo() echocontract.EchoServiceGateways {
	return p.echoServiceGateways
}
