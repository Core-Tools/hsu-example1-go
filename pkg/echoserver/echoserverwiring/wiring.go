package echoserverwiring

import (
	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduleapi"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/modulewiring"
	"github.com/core-tools/hsu-echo/pkg/echoapi"
	"github.com/core-tools/hsu-echo/pkg/echocontract"
	"github.com/core-tools/hsu-echo/pkg/echoserver/echoserverdomain"
)

func init() {
	// Self-register in global registry
	moduleDesc := modulewiring.ModuleDescriptor[
		echoserverdomain.EchoServiceProvider,
		echocontract.EchoServiceGateways,
		echocontract.EchoServiceHandlers,
	]{
		ServiceProviderFactoryFunc: NewEchoServiceProvider,
		ModuleFactoryFunc:          echoserverdomain.NewEchoModule,
		HandlersRegistrarFunc:      echoapi.EchoHandlersRegistrar,
		DirectClosureEnablerFunc:   echoapi.EchoDirectClosureEnabler,
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
