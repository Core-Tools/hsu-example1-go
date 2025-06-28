package domain

import (
	"context"
)

type Contract interface {
	Echo(ctx context.Context, message string) (string, error)
}
