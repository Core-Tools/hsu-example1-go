package contract

import (
	"context"
)

type Contract1Echo1Func func(ctx context.Context, message string) (string, error)
type Contract1Echo2Func func(ctx context.Context, message string) (string, error)

type Contract1 interface {
	Echo1(ctx context.Context, message string) (string, error)
	Echo2(ctx context.Context, message string) (string, error)
}
