package domain

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-echo/pkg/contract"
)

func NewService1(logger logging.Logger) contract.Contract1 {
	return &service1{
		logger: logger,
	}
}

type service1 struct {
	logger logging.Logger
}

func (s *service1) Echo1(ctx context.Context, message string) (string, error) {
	s.logger.Debugf("Module: Service1.Echo1 called: %s", message)
	return message, nil
}

func (s *service1) Echo2(ctx context.Context, message string) (string, error) {
	s.logger.Debugf("Module: Service1.Echo2 called: %s", message)
	return message, nil
}
