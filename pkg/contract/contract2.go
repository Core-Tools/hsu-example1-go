package contract

import (
	"context"
)

type Contract2Echo1Func func(ctx context.Context, message string) (string, error)
type Contract2Echo2Func func(ctx context.Context, message string) (string, error)

type Contract2 interface {
	Echo1(ctx context.Context, message string) (string, error)
	Echo2(ctx context.Context, message string) (string, error)
}
