package main

import (
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

// n == 40, workerNum == 6 -> duration == 664ms
// on MacBook Pro 15"" 2019, 16 GB, 6-Core i7
func TestFibonacci(t *testing.T) {
	provider := &Provider{}
	startTime := time.Now()
	result, err := provider.Fibonacci(50, 6)
	assert.Nil(t, err)
	duration := time.Now().Sub(startTime)
	t.Logf("duration: %v", duration)
	t.Logf("result: %d", result)
}
