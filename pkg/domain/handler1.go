package echodomain

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	echocontract "github.com/core-tools/hsu-echo/pkg/api/contract"
)

func NewHandler1(logger logging.Logger) echocontract.Service1 {
	return &handler1{
		logger: logger,
	}
}

type handler1 struct {
	logger logging.Logger
}

func (h *handler1) Echo1(ctx context.Context, message string) (string, error) {
	h.logger.Debugf("Module: Service1.Echo1 called: %s", message)
	return message, nil
}

func (h *handler1) Echo2(ctx context.Context, message string) (string, error) {
	h.logger.Debugf("Module: Service1.Echo2 called: %s", message)
	return message, nil
}
