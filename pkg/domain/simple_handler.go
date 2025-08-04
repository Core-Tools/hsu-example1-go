package domain

import (
	"context"

	"github.com/core-tools/hsu-core/pkg/logging"
	"github.com/core-tools/hsu-echo/pkg/contract"
)

func NewSimpleHandler(logger logging.Logger) contract.Contract {
	return &simpleHandler{
		logger: logger,
	}
}

type simpleHandler struct {
	logger logging.Logger
}

func (h *simpleHandler) Echo(ctx context.Context, message string) (string, error) {
	return "go-simple-echo: " + message, nil
}
