package main

import (
	_ "net/http/pprof"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/impl"
	_ "github.com/apache/dubbo-go/protocol/jsonrpc"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

import (
	"github.com/dubbogo/dubbo-go-benchmark/common"
)

var (
	survivalTimeout = int(3e9)
)

// they are necessary:
// 		export CONF_PROVIDER_FILE_PATH="xxx"
// 		export APP_LOG_CONF_FILE="xxx"

func main() {
	config.Load()

	common.InitProfiling("7070")

	common.InitSignal(survivalTimeout)
}
