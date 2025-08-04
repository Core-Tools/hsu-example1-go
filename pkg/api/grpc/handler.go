package grpcapi

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-echo/pkg/contract"
	"github.com/core-tools/hsu-echo/pkg/generated/api/proto"

	"google.golang.org/grpc"
)

func RegisterGRPCHandler(grpcServerRegistrar grpc.ServiceRegistrar, handler interface{}, logger logging.Logger) {
	proto.RegisterEchoServiceServer(grpcServerRegistrar, &grpcServerHandler{
		handler: handler.(contract.Contract),
		logger:  logger,
	})
}

type grpcServerHandler struct {
	proto.UnimplementedEchoServiceServer
	handler contract.Contract
	logger  logging.Logger
}

func (h *grpcServerHandler) Echo(ctx context.Context, echoRequest *proto.EchoRequest) (*proto.EchoResponse, error) {
	response, err := h.handler.Echo(ctx, echoRequest.Message)
	if err != nil {
		h.logger.Errorf("Echo server handler: %v", err)
		return nil, err
	}
	h.logger.Debugf("Echo server handler done")
	return &proto.EchoResponse{Message: response}, nil
}
