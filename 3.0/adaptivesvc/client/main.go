package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

import (
	clusterutils "dubbo.apache.org/dubbo-go/v3/cluster/utils"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	testerpkg "github.com/dubbogo/tools/pkg/tester"
)

const (
	Tps      = "TPS"
	Parallel = "PARALLEL"
	// Duration should be a string representing a time,
	// like "1h", "30m", etc.
	Duration = "DURATION"
	FuncName = "FUNC_NAME"

	// Supported FuncNames
	Fibonacci = "FIBONACCI"
	Sleep     = "SLEEP"

	// FuncName == "FIBONACCI"
	FibonacciN = "FIBONACCI_N"

	// FuncName == "SLEEP"
	SleepDuration = "SLEEP_DURATION"
)

type Provider struct {
	Fibonacci func(ctx context.Context, n, workerNum int64) (int64, error)
	Sleep     func(ctx context.Context, duration int64) (int64, error)
}

func main() {
	provider := &Provider{}
	config.SetConsumerService(provider)

	if err := config.Load(); err != nil {
		panic(err)
	}

	var (
		tps, parallel      int
		duration, funcName string
		err                error
	)

	ctx := context.TODO()
	if tps, err = strconv.Atoi(os.Getenv(Tps)); err != nil {
		panic(fmt.Errorf("env %s is required: %w", Tps, err))
	}
	logger.Infof("TPS is set to %d.", tps)
	if parallel, err = strconv.Atoi(os.Getenv(Parallel)); err != nil {
		panic(fmt.Errorf("env %s is required: %w", Parallel, err))
	}
	logger.Infof("Parallel is set to %d.", parallel)
	if duration = os.Getenv(Duration); duration == "" {
		panic(fmt.Errorf("%s is required", Duration))
	}
	if funcName = os.Getenv(FuncName); funcName == "" {
		panic(fmt.Errorf("%s is required", FuncName))
	}

	doInvoke := func(uid int) {
		switch funcName {
		case Fibonacci:
			if result, err := fibonacci(ctx, provider); err != nil {
				if clusterutils.DoesAdaptiveServiceReachLimitation(err) {
					logger.Infof("Reach Limitation")
				} else {
					panic(err)
				}
			} else {
				fmt.Printf("%s result: %d\n", Fibonacci, result)
			}
		case Sleep:
			sleep(ctx, provider)
			fmt.Printf("sleep task was finished")
		default:
			panic(fmt.Sprintf("%s is an unsupported function", funcName))
		}
	}

	tester := testerpkg.NewStressTester()
	tester.
		SetTPS(tps).
		SetDuration(duration).
		SetTestFn(doInvoke).
		SetUserNum(parallel).
		Run()

	fmt.Printf("Sent request num: %d", tester.GetTransactionNum())
	fmt.Printf("TPS: %.2f\n", tester.GetTPS())
	fmt.Printf("RT: %.2fs\n", tester.GetAverageRTSeconds())
}

func fibonacci(ctx context.Context, provider *Provider) (result int64, err error) {
	var n int
	if n, err = strconv.Atoi(os.Getenv(FibonacciN)); err != nil {
		return
	}

	result, err = provider.Fibonacci(ctx, int64(n), 1)
	return
}

func sleep(ctx context.Context, provider *Provider) {
	var (
		duration int
		err      error
	)
	if duration, err = strconv.Atoi(os.Getenv(SleepDuration)); err != nil {
		panic(err)
	}
	_, _ = provider.Sleep(ctx, int64(duration))
}
