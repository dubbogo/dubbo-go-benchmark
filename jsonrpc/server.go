package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

import (
	"github.com/AlexStocks/goext/time"
	log "github.com/AlexStocks/log4go"
)

import (
	_ "github.com/dubbo/go-for-apache-dubbo/cluster/cluster_impl"
	_ "github.com/dubbo/go-for-apache-dubbo/cluster/loadbalance"
	_ "github.com/dubbo/go-for-apache-dubbo/common/proxy/proxy_factory"
	"github.com/dubbo/go-for-apache-dubbo/config"
	_ "github.com/dubbo/go-for-apache-dubbo/filter/impl"
	_ "github.com/dubbo/go-for-apache-dubbo/protocol/jsonrpc"
	_ "github.com/dubbo/go-for-apache-dubbo/registry/protocol"
	_ "github.com/dubbo/go-for-apache-dubbo/registry/zookeeper"
)

var (
	survivalTimeout = int(3e9)
)

// they are necessary:
// 		export CONF_PROVIDER_FILE_PATH="xxx"
// 		export APP_LOG_CONF_FILE="xxx"

func main() {
	_, proMap := config.Load()
	if proMap == nil {
		panic("proMap is nil")
	}
	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		log.Info("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
		// reload()
		default:
			go gxtime.Future(survivalTimeout, func() {
				log.Warn("app exit now by force...")
				os.Exit(1)
			})

			// 要么fastFailTimeout时间内执行完毕下面的逻辑然后程序退出，要么执行上面的超时函数程序强行退出
			fmt.Println("provider app exit now...")
			return
		}
	}
}
