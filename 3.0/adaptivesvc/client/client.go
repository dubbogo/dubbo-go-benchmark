package main

import "dubbo.apache.org/dubbo-go/v3/config"

func main() {
	provider := &Provider{}
	config.SetConsumerService(provider)

	if err := config.Load(config.WithPath("./dubbogo.yml")); err != nil {
		panic(err)
	}
	select {}
}
