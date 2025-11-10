package echogrpcapi

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-echo/pkg/echocontract"
	"github.com/core-tools/hsu-echo/pkg/generated/api/proto"

	"google.golang.org/grpc"
)

func RegisterGRPCHandler1(grpcServerRegistrar grpc.ServiceRegistrar, serviceHandler echocontract.Service1, logger logging.Logger) error {
	proto.RegisterEchoServiceServer(grpcServerRegistrar, &grpcServiceHandler1{
		serviceHandler: serviceHandler,
		logger:         logger,
	})
	return nil
}

type grpcServiceHandler1 struct {
	proto.UnimplementedEchoServiceServer
	serviceHandler echocontract.Service1
	logger         logging.Logger
}

func (h *grpcServiceHandler1) Echo(ctx context.Context, echoRequest *proto.EchoRequest) (*proto.EchoResponse, error) {
	h.logger.Debugf("gRPC Handler: Service1.Echo1 request: %s", echoRequest.Message)
	response, err := h.serviceHandler.Echo1(ctx, echoRequest.Message)
	if err != nil {
		h.logger.Errorf("gRPC Handler: Service1.Echo1 request failed: %v", err)
		return nil, err
	}
	h.logger.Debugf("gRPC Handler: Service1.Echo1 response: %s", response)
	return &proto.EchoResponse{Message: response}, nil
}
