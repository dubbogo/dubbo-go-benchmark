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
	"context"
	"os"
	"strconv"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/tools/pkg/stressTest"
)

import (
	"github.com/dubbogo/dubbo-go-benchmark/3.0/triple/api"
)

var grpcGreeterImpl = new(api.GreeterClientImpl)

func main() {
	config.SetConsumerService(grpcGreeterImpl)
	if err := config.Load(config.WithPath("./dubbogo.yml")); err != nil {
		panic(err)
	}

	ctx := context.Background()
	tpsNum, _ := strconv.Atoi(os.Getenv("tps"))
	parallel, _ := strconv.Atoi(os.Getenv("parallel"))
	payloadLen, _ := strconv.Atoi(os.Getenv("payload"))
	req := &api.HelloRequest{
		Name: "laurence" + string(make([]byte, payloadLen)),
	}
	stressTest.NewStressTestConfigBuilder().SetTPS(tpsNum).SetDuration("1h").SetParallel(parallel).Build().Start(func() {
		if _, err := grpcGreeterImpl.SayHello(ctx, req); err != nil {
			panic(err)
		}
	})
}
