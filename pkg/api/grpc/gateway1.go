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

func NewGRPCGateway1(clientConnection moduleproto.ProtocolClientConnection, logger logging.Logger) moduletypes.ServiceGateway {
	grpcClientConnection := clientConnection.(*grpc.ClientConn)
	grpcClient := proto.NewEchoServiceClient(grpcClientConnection)
	return &grpcGateway1{
		grpcClient: grpcClient,
		logger:     logger,
	}
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

var _ contract.Contract1 = (*grpcGateway1)(nil)
