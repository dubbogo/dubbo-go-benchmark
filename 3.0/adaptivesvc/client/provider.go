package main

import (
	"context"
)

type Provider struct {
	Fibonacci func(ctx context.Context, n, workerNum int64) (int64, error)
}
