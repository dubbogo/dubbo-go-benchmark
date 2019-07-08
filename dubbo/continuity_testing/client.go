package main

import (
	"context"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"time"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/impl"
	"github.com/dubbogo/hessian2"

	_ "github.com/apache/dubbo-go/registry/zookeeper"

	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
)

import (
	"github.com/dubbogo/go-for-apache-dubbo-benchmark/common"
)

// they are necessary:
// 		export CONF_CONSUMER_FILE_PATH="xxx"
// 		export APP_LOG_CONF_FILE="xxx"

var (
	arg = flag.Int("r", 1, "size of arg(1 = 300B)")
	ti  = flag.Int("t", 30, "interval time(ms)")

	Arg             = common.GetString(*arg)
	survivalTimeout = int(3e9)
)

func main() {
	common.InitProfiling("7072")
	flag.Parse()

	hessian.RegisterJavaEnum(Gender(MAN))
	hessian.RegisterJavaEnum(Gender(WOMAN))
	hessian.RegisterPOJO(&User{})

	config.Load()

	time.Sleep(3e9)

	for {
		for i := 0; i < 10; i++ {
			go func() {
				user := &User{}
				err := userProvider.GetUser(context.TODO(), []interface{}{"A003", Arg}, user)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}
		<-time.After(time.Duration(int64(*ti) * int64(time.Millisecond)))
	}

	common.InitSignal(survivalTimeout)
}
