package main

import (
	_ "net/http/pprof"
)

import (
	"github.com/apache/dubbo-go/config"
	"github.com/dubbogo/hessian2"

	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"

	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/filter/impl"

	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"

	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

import (
	"github.com/dubbogo/dubbo-go-benchmark/common"
)

var (
	survivalTimeout = int(3e9)
)

func main() {

	// ------for hessian2------
	hessian.RegisterJavaEnum(Gender(MAN))
	hessian.RegisterJavaEnum(Gender(WOMAN))
	hessian.RegisterPOJO(&User{})
	// ------------

	config.Load()

	common.InitProfiling("7070")

	common.InitSignal(survivalTimeout)
}
