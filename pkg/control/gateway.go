package control

import (
	"context"

	"github.com/core-tools/hsu-echo/pkg/domain"
	"github.com/core-tools/hsu-echo/pkg/generated/api/proto"
	"github.com/core-tools/hsu-echo/pkg/logging"

	"google.golang.org/grpc"
)

func NewGRPCClientGateway(grpcClientConnection grpc.ClientConnInterface, logger logging.Logger) domain.Contract {
	grpcClient := proto.NewEchoServiceClient(grpcClientConnection)
	return &grpcClientGateway{
		grpcClient: grpcClient,
		logger:     logger,
	}
}

type grpcClientGateway struct {
	grpcClient proto.EchoServiceClient
	logger     logging.Logger
}

func (gw *grpcClientGateway) Echo(ctx context.Context, message string) (string, error) {
	response, err := gw.grpcClient.Echo(ctx, &proto.EchoRequest{Message: message})
	if err != nil {
		gw.logger.Errorf("Echo client gateway: %v", err)
		return "", err
	}
	gw.logger.Debugf("Echo client gateway done")
	return response.Message, nil
}
