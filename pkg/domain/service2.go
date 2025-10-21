package domain

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-echo/pkg/contract"
)

func NewService2(logger logging.Logger) contract.Contract2 {
	return &service2{
		logger: logger,
	}
}

type service2 struct {
	logger logging.Logger
}

func (s *service2) Echo1(ctx context.Context, message string) (string, error) {
	s.logger.Debugf("Module: Service2.Echo1 called: %s", message)
	return "go-service2-echo1: " + message, nil
}

func (s *service2) Echo2(ctx context.Context, message string) (string, error) {
	s.logger.Debugf("Module: Service2.Echo2 called: %s", message)
	return "go-service2-echo2: " + message, nil
}
