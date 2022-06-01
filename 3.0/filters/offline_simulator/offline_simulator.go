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

package offline_simulator

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
)

type ServerState string

const (
	ServerStateOnline  ServerState = "Online"
	ServerStateOffline ServerState = "Offline"

	// MinOnlineDuration should be a string representing a time,
	// like "1h", "30m", etc.
	MinOnlineDuration  = "MIN_ONLINE_DURATION"
	MaxOnlineDuration  = "MAX_ONLINE_DURATION"
	MinOfflineDuration = "MIN_OFFLINE_DURATION"
	MaxOfflineDuration = "MAX_OFFLINE_DURATION"
	//OfflineRatio should be a decimal between 0 and 1.
	OfflineRatio = "OFFLINE_RATIO"
)

var ErrServerOffline = fmt.Errorf("server is offline")

type OfflineSimulator struct {
	State              ServerState
	OfflineRatio       float64
	MinOnlineDuration  time.Duration  //default is about 1s
	MaxOnlineDuration  *time.Duration //optional
	MinOfflineDuration time.Duration  //default is about 1s
	MaxOfflineDuration *time.Duration //optional
	LastTransferTime   time.Time
}

func init() {
	extension.SetFilter("offlineSimulator", NewOfflineSimulator)
}

func NewOfflineSimulator() filter.Filter {
	var (
		err                error
		offlineRatio       float64
		minOnlineDuration  time.Duration
		maxOnlineDuration  *time.Duration
		minOfflineDuration time.Duration
		maxOfflineDuration *time.Duration
	)

	if offlineRatioStr := os.Getenv(OfflineRatio); offlineRatioStr != "" {
		offlineRatio, err = strconv.ParseFloat(offlineRatioStr, 64)
		if err != nil {
			panic(fmt.Errorf("%s should be a decimal", OfflineRatio))
		}
		if offlineRatio > 1 || offlineRatio < 0 {
			panic(fmt.Errorf("%s should be a decimal between 0 and 1 ", OfflineRatio))
		}
	}

	if minOnlineDurationStr := os.Getenv(MinOnlineDuration); minOnlineDurationStr != "" {
		minOnlineDuration, err = time.ParseDuration(minOnlineDurationStr)
		if err != nil {
			panic(fmt.Errorf("%s should be a string representing a time, like \"1h\", \"30m\", etc", MinOnlineDuration))
		}
	}

	if maxOnlineDurationStr := os.Getenv(MaxOnlineDuration); maxOnlineDurationStr != "" {
		duration, err := time.ParseDuration(maxOnlineDurationStr)
		if err != nil {
			panic(fmt.Errorf("%s should be a string representing a time, like \"1h\", \"30m\", etc", MaxOnlineDuration))
		}
		maxOnlineDuration = &duration
	}

	if minOfflineDurationStr := os.Getenv(MinOfflineDuration); minOfflineDurationStr != "" {
		minOfflineDuration, err = time.ParseDuration(minOfflineDurationStr)
		if err != nil {
			panic(fmt.Errorf("%s should be a string representing a time, like \"1h\", \"30m\", etc", MinOfflineDuration))
		}
	}

	if maxOfflineDurationStr := os.Getenv(MaxOfflineDuration); maxOfflineDurationStr != "" {
		duration, err := time.ParseDuration(maxOfflineDurationStr)
		if err != nil {
			panic(fmt.Errorf("%s should be a string representing a time, like \"1h\", \"30m\", etc", MaxOfflineDuration))
		}
		maxOfflineDuration = &duration
	}

	s := &OfflineSimulator{
		State:              ServerStateOnline,
		OfflineRatio:       offlineRatio,
		MinOnlineDuration:  minOnlineDuration,
		MaxOnlineDuration:  maxOnlineDuration,
		MinOfflineDuration: minOfflineDuration,
		MaxOfflineDuration: maxOfflineDuration,
		LastTransferTime:   time.Now(),
	}
	go s.Run()
	return s
}

//Run an offline simulator
//1. if the duration of the current state is less than the minimum limit, sleep until greater than it;
//2. if the duration is greater than the maximum limit (if set), switch the state immediately;
//3. otherwise, switch the state with a certain probability every second.
func (f *OfflineSimulator) Run() {
	if f.OfflineRatio <= 0 {
		return
	}
	for {
		now := time.Now()
		switch f.State {
		case ServerStateOnline:
			if duration := now.Sub(f.LastTransferTime); duration < f.MinOnlineDuration {
				time.Sleep(f.MinOnlineDuration - duration + time.Second)
			} else if (f.MaxOnlineDuration != nil && duration > *f.MaxOnlineDuration) || rand.Float64() < f.OfflineRatio {
				f.State = ServerStateOffline
				f.LastTransferTime = time.Now()
			} else {
				time.Sleep(time.Second)
			}
		case ServerStateOffline:
			if duration := now.Sub(f.LastTransferTime); duration < f.MinOfflineDuration {
				time.Sleep(f.MinOfflineDuration - duration + time.Second)
			} else if (f.MaxOfflineDuration != nil && duration > *f.MaxOfflineDuration) || rand.Float64() > f.OfflineRatio {
				f.State = ServerStateOnline
				f.LastTransferTime = time.Now()
			} else {
				time.Sleep(time.Second)
			}
		}
	}
}

func IsServerOfflineErr(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == ErrServerOffline.Error()
}

func (f *OfflineSimulator) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	if f.State == ServerStateOffline {
		return &protocol.RPCResult{
			Attrs: nil,
			Err:   ErrServerOffline,
			Rest:  nil,
		}
	}
	return invoker.Invoke(ctx, invocation)
}

func (f *OfflineSimulator) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	if f.State == ServerStateOffline {
		return &protocol.RPCResult{
			Attrs: nil,
			Err:   ErrServerOffline,
			Rest:  nil,
		}
	}
	return result
}
