package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"github.com/dubbogo/tools/pkg/stressTest"
)

const (
	Tps = "TPS"
	Parallel = "PARALLEL"
	// Duration should be a string representing a time,
	// like "1h", "30m", etc.
	Duration = "DURATION"
	FuncName = "FUNC_NAME"

	// Supported FuncNames
	Fibonacci = "FIBONACCI"

	// FuncName == "FIBONACCI"
	FibonacciN = "FIBONACCI_N"
	FibonacciWorkerNum = "FIBONACCI_WORKER_NUM"
)

func main() {
	provider := &Provider{}
	config.SetConsumerService(provider)

	if err := config.Load(config.WithPath("./dubbogo.yml")); err != nil {
		panic(err)
	}

	var (
		tps, parallel int
		duration, funcName string
		err error
	)

	ctx := context.TODO()
	if tps, err = strconv.Atoi(os.Getenv(Tps)); err != nil {
		panic(err)
	}
	if parallel, err = strconv.Atoi(os.Getenv(Parallel)); err != nil {
		panic(err)
	}
	if duration = os.Getenv(Duration); duration == "" {
		panic(fmt.Errorf("%s is required", Duration))
	}
	if funcName = os.Getenv(FuncName); funcName == "" {
		panic(fmt.Errorf("%s is required", FuncName))
	}

	doInvoke := func() {
		switch funcName {
		case Fibonacci:
			if result, err := fibonacci(ctx, provider); err != nil {
				panic(err)
			} else {
				fmt.Printf("%s result: %d\n", Fibonacci, result)
			}
		default:
			panic(fmt.Sprintf("%s is an unsupported function", funcName))
		}
	}

	stressTest.NewStressTestConfigBuilder().
		SetTPS(tps).
		SetDuration(duration).
		SetParallel(parallel).Build().
		Start(doInvoke)
}

func fibonacci(ctx context.Context, provider *Provider) (result int, err error) {
	var (
		n, workNum int
	)
	if n, err = strconv.Atoi(os.Getenv(FibonacciN)); err != nil {
		panic(err)
	}
	if workNum, err = strconv.Atoi(os.Getenv(FibonacciWorkerNum)); err != nil {
		panic(err)
	}

	result, err = provider.Fibonacci(ctx, n, workNum)
	return
}
