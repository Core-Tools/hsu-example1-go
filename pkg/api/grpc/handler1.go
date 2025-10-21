package grpcapi

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduleproto"
	"github.com/core-tools/hsu-core/pkg/modulemanagement/moduletypes"
	"github.com/core-tools/hsu-echo/pkg/contract"
	"github.com/core-tools/hsu-echo/pkg/generated/api/proto"

	"google.golang.org/grpc"
)

func RegisterGRPCHandler1(grpcServerRegistrar grpc.ServiceRegistrar, service contract.Contract1, logger logging.Logger) {
	proto.RegisterEchoServiceServer(grpcServerRegistrar, &grpcServerHandler1{
		service: service,
		logger:  logger,
	})
}

type grpcServerHandler1 struct {
	proto.UnimplementedEchoServiceServer
	service contract.Contract1
	logger  logging.Logger
}

func (h *grpcServerHandler1) Echo(ctx context.Context, echoRequest *proto.EchoRequest) (*proto.EchoResponse, error) {
	h.logger.Debugf("gRPC Handler: Service1.Echo1 request: %s", echoRequest.Message)
	response, err := h.service.Echo1(ctx, echoRequest.Message)
	if err != nil {
		h.logger.Errorf("gRPC Handler: Service1.Echo1 request failed: %v", err)
		return nil, err
	}
	h.logger.Debugf("gRPC Handler: Service1.Echo1 response: %s", response)
	return &proto.EchoResponse{Message: response}, nil
}

func RegisterHandlers(serviceHandlersMap moduletypes.ServiceHandlersMap, protocolServerRegistrar moduleproto.ProtocolServerRegistrar, logger logging.Logger) {
	grpcServerRegistrar := protocolServerRegistrar.(grpc.ServiceRegistrar)

	service1ID := moduletypes.ServiceID("service1")
	abstractService1 := serviceHandlersMap[service1ID]
	typifiedService1 := abstractService1.(contract.Contract1)
	RegisterGRPCHandler1(grpcServerRegistrar, typifiedService1, logger)
}
