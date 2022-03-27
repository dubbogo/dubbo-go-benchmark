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
	"sync"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

var ErrNShouldGreaterThanZero = fmt.Errorf("n should greater than zero")

type Provider struct{}

func (*Provider) Fibonacci(n, workerNum int64) (int64, error) {
	var (
		result int64
		err    error
		wg     sync.WaitGroup
	)
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
	config.SetProviderService(&Provider{})
	if err := config.Load(); err != nil {
		panic(err)
	}
	select {}
}
