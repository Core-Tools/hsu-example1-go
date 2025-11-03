package echodomain

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	echocontract "github.com/core-tools/hsu-echo/pkg/api/contract"
)

func NewHandler2(logger logging.Logger) echocontract.Service2 {
	return &handler2{
		logger: logger,
	}
}

type handler2 struct {
	logger logging.Logger
}

func (h *handler2) Echo1(ctx context.Context, message string) (string, error) {
	h.logger.Debugf("Module: Service2.Echo1 called: %s", message)
	return "go-service2-echo1: " + message, nil
}

func (h *handler2) Echo2(ctx context.Context, message string) (string, error) {
	h.logger.Debugf("Module: Service2.Echo2 called: %s", message)
	return "go-service2-echo2: " + message, nil
}
