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
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)


var (
	userProvider = &UserProvider{}
)


type UserProviderProxy struct {
	userProvider *UserProvider
}

func (u *UserProviderProxy) GetUser(ctx context.Context, userID *Request) (*User, error) {
	return u.userProvider.GetUser(ctx, userID)
	//return &User{
	//	ID: "12345",
	//	Name: "Hello" + userID,
	//	Age: 21,
	//}, nil
}


// need to setup environment variable "DUBBO_GO_CONFIG_PATH" to "conf/dubbogo.yml" before run
func main() {
	config.SetConsumerService(userProvider)
	config.SetProviderService(&UserProviderProxy{
		userProvider: userProvider,
	})

	err := config.Load(config.WithPath("./dubbogo.yml"))
	if err != nil {
		panic(err)
	}
	select {

	}
}
