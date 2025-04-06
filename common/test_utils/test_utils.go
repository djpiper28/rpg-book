package testutils

import (
	"context"
	"time"
)

func NewTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Minute/2)
}
