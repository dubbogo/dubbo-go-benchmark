package main

import "context"

type Provider struct {
	Fibonacci func(ctx context.Context, n, worker int) (int, error)
}
