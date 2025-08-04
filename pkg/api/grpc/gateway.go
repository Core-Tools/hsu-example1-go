package grpcapi

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-echo/pkg/generated/api/proto"

	"google.golang.org/grpc"
)

func NewGRPCGateway(grpcClientConnection grpc.ClientConnInterface, logger logging.Logger) interface{} {
	grpcClient := proto.NewEchoServiceClient(grpcClientConnection)
	return &grpcGateway{
		grpcClient: grpcClient,
		logger:     logger,
	}
}

type grpcGateway struct {
	grpcClient proto.EchoServiceClient
	logger     logging.Logger
}

func (gw *grpcGateway) Echo(ctx context.Context, message string) (string, error) {
	response, err := gw.grpcClient.Echo(ctx, &proto.EchoRequest{Message: message})
	if err != nil {
		gw.logger.Errorf("Echo client gateway: %v", err)
		return "", err
	}
	gw.logger.Debugf("Echo client gateway done")
	return response.Message, nil
}
