package echoapp

import (
	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduleapi"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/modulewiring"
	echoapi "github.com/core-tools/hsu-echo/pkg/api"
	echocontract "github.com/core-tools/hsu-echo/pkg/api/contract"
	echodomain "github.com/core-tools/hsu-echo/pkg/domain"
)

func init() {
	// Self-register in global registry
	moduleDesc := modulewiring.ModuleDescriptor[
		echodomain.EchoServiceProvider,
		echocontract.EchoServiceGateways,
		echocontract.EchoServiceHandlers,
	]{
		ServiceProviderFactoryFunc:   NewEchoServiceProvider,
		ModuleFactoryFunc:            echodomain.NewEchoModule,
		HandlersRegistrarFactoryFunc: echoapi.NewEchoHandlersRegistrar,
		DirectClosureEnableFunc:      echoapi.EchoDirectClosureEnable,
	}
	modulewiring.RegisterModule("echo", moduleDesc)
}

func NewEchoServiceProvider(serviceConnector moduleapi.ServiceConnector, logger logging.Logger) moduletypes.ServiceProviderHandle {
	serviceProvider := &echoServiceProvider{}
	return moduletypes.ServiceProviderHandle{
		ServiceProvider: serviceProvider,
	}
}

type echoServiceProvider struct{}
