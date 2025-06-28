package control

import (
	"context"

	"github.com/core-tools/hsu-echo/pkg/domain"
	"github.com/core-tools/hsu-echo/pkg/generated/api/proto"
	"github.com/core-tools/hsu-echo/pkg/logging"

	"google.golang.org/grpc"
)

func RegisterGRPCServerHandler(grpcServerRegistrar grpc.ServiceRegistrar, handler domain.Contract, logger logging.Logger) {
	proto.RegisterEchoServiceServer(grpcServerRegistrar, &grpcServerHandler{
		handler: handler,
		logger:  logger,
	})
}

type grpcServerHandler struct {
	proto.UnimplementedEchoServiceServer
	handler domain.Contract
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
