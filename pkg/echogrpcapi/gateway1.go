package echogrpcapi

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-echo/pkg/echocontract"
	"github.com/core-tools/hsu-echo/pkg/generated/api/proto"

	"google.golang.org/grpc"
)

func NewGRPCGateway1(grpcClientConnection *grpc.ClientConn, logger logging.Logger) (echocontract.Service1, error) {
	return &grpcGateway1{
		grpcClient: proto.NewEchoServiceClient(grpcClientConnection),
		logger:     logger,
	}, nil
}

type grpcGateway1 struct {
	grpcClient proto.EchoServiceClient
	logger     logging.Logger
}

func (gw *grpcGateway1) Echo1(ctx context.Context, message string) (string, error) {
	gw.logger.Debugf("gRPC Gateway: Service1.Echo1 request: %s", message)
	response, err := gw.grpcClient.Echo(ctx, &proto.EchoRequest{Message: message})
	if err != nil {
		gw.logger.Errorf("gRPC Gateway: Service1.Echo1 request failed: %v", err)
		return "", err
	}
	gw.logger.Debugf("gRPC Gateway: Service1.Echo1 response: %s", response.Message)
	return response.Message, nil
}

func (gw *grpcGateway1) Echo2(ctx context.Context, message string) (string, error) {
	gw.logger.Debugf("gRPC Gateway: Service1.Echo2 request: %s", message)
	response, err := gw.grpcClient.Echo(ctx, &proto.EchoRequest{Message: message})
	if err != nil {
		gw.logger.Errorf("gRPC Gateway: Service1.Echo2 request failed: %v", err)
		return "", err
	}
	gw.logger.Debugf("gRPC Gateway: Service1.Echo2 response: %s", response.Message)
	return response.Message, nil
}

var _ echocontract.Service1 = (*grpcGateway1)(nil)
