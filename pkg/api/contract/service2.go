package echocontract

import (
	"context"
)

type Service2 interface {
	Echo1(ctx context.Context, message string) (string, error)
	Echo2(ctx context.Context, message string) (string, error)
}
