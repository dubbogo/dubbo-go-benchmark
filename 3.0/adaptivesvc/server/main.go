/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	_ "github.com/dubbogo/dubbo-go-benchmark/3.0/filters/metrics_collector"
	_ "github.com/dubbogo/dubbo-go-benchmark/3.0/filters/offline_simulator"
)

const (
	// TimeoutDuration should be a string representing a time,
	// like "1h", "30m", etc.
	TimeoutDuration = "TIMEOUT_DURATION"
	//TimeoutRatio should be a decimal between 0 and 1.
	TimeoutRatio = "TIMEOUT_RATIO"
	RandSeed     = "RAND_SEED"
)

var (
	timeoutRatio    float64
	timeoutDuration time.Duration
)

var ErrNShouldGreaterThanZero = fmt.Errorf("n should greater than zero")

type Provider struct{}

func (*Provider) Fibonacci(n, workerNum int64) (int64, error) {
	var (
		result int64
		err    error
		wg     sync.WaitGroup
	)

	if rand.Float64() < timeoutRatio {
		time.Sleep(timeoutDuration)
	}

	for i := 0; i < int(workerNum); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ret, e := fibonacci(n); e != nil {
				err = e
			} else {
				result = ret
			}
		}()
	}
	wg.Wait()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func fibonacci(n int64) (int64, error) {
	if n < 0 {
		return 0, ErrNShouldGreaterThanZero
	}
	if n < 2 {
		return n, nil
	}

	f1, err := fibonacci(n - 1)
	if err != nil {
		return 0, err
	}
	f2, err := fibonacci(n - 2)
	if err != nil {
		return 0, err
	}

	return f1 + f2, nil
}

func (*Provider) Sleep(duration int64) (int64, error) {
	time.Sleep(time.Duration(duration) * time.Millisecond)
	return 0, nil
}

func main() {

	var (
		err      error
		randSeed int64
	)
	if randSeedStr := os.Getenv(RandSeed); randSeedStr != "" {
		randSeed, err = strconv.ParseInt(randSeedStr, 10, 64)
		if err != nil {
			panic(fmt.Errorf("%s should be a integer", RandSeed))
		}
	} else {
		randSeed = time.Now().Unix()
	}
	rand.Seed(randSeed)

	if timeoutRateStr := os.Getenv(TimeoutRatio); timeoutRateStr != "" {
		timeoutRatio, err = strconv.ParseFloat(timeoutRateStr, 64)
		if err != nil {
			panic(fmt.Errorf("%s should be a decimal", TimeoutRatio))
		}
		if timeoutRatio > 1 || timeoutRatio < 0 {
			panic(fmt.Errorf("%s should be a decimal between 0 and 1 ", TimeoutRatio))
		}
	}

	if timeoutDurationStr := os.Getenv(TimeoutDuration); timeoutDurationStr == "" && timeoutRatio > 0 {
		panic(fmt.Errorf("%s is required", TimeoutDuration))
	} else if timeoutDurationStr != "" {
		timeoutDuration, err = time.ParseDuration(timeoutDurationStr)
		if err != nil {
			panic(fmt.Errorf("%s should be a string representing a time, like \"1h\", \"30m\", etc", TimeoutDuration))
		}
	}

	config.SetProviderService(&Provider{})
	if err := config.Load(); err != nil {
		panic(err)
	}

	select {}
}
