package domain

import (
	"context"

	"github.com/core-tools/hsu-echo/pkg/domain"
	"github.com/core-tools/hsu-echo/pkg/logging"
)

func NewSimpleHandler(logger logging.Logger) domain.Contract {
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
