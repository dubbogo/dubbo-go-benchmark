package main

import (
	"context"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"sync"
	"sync/atomic"
	"time"
)

import (
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/protocol/jsonrpc"
	_ "github.com/apache/dubbo-go/registry/protocol"
	"github.com/montanaflynn/stats"

	_ "github.com/apache/dubbo-go/filter/impl"

	_ "github.com/apache/dubbo-go/cluster"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"

	_ "github.com/apache/dubbo-go/cluster/cluster_impl"

	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

import (
	"github.com/dubbogo/dubbo-go-benchmark/common"
)

// they are necessary:
// 		export CONF_CONSUMER_FILE_PATH="xxx"
// 		export APP_LOG_CONF_FILE="xxx"

var (
	concurrency = flag.Int("c", 1, "concurrency")
	total       = flag.Int("n", 1, "total requests for all clients")
	arg         = flag.Int("r", 2, "size of arg(1 = 300B)")

	Arg             = common.GetString(*arg)
	survivalTimeout = int(3e9)
)

func main() {
	common.InitProfiling("7071")
	flag.Parse()
	Arg = common.GetString(*arg)
	conc, tn, err := common.CheckArgs(*concurrency, *total)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	n := conc
	m := tn / n

	fmt.Printf("concurrency: %d\nrequests per client: %d\n\n", n, m)

	var wg sync.WaitGroup
	wg.Add(n * m)

	fmt.Printf("sent total %d messages, %d message per client", n*m, m)

	config.Load()

	time.Sleep(3e9)

	var trans uint64
	var transOK uint64

	d := make([][]int64, n, n)

	totalT := time.Now().UnixNano()
	for i := 0; i < n; i++ {
		dt := make([]int64, 0, m)
		d = append(d, dt)

		go func(i int) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered in %v", r)
				}
			}()

			for j := 0; j < m; j++ {
				t := time.Now().UnixNano()
				user := &User{}
				err := userProvider.GetUser(context.TODO(), []interface{}{"A003", Arg}, user)

				t = time.Now().UnixNano() - t

				d[i] = append(d[i], t)

				if err == nil && user.Id != "" {
					atomic.AddUint64(&transOK, 1)
				}

				if err != nil {
					fmt.Printf(err.Error())
				}

				atomic.AddUint64(&trans, 1)
				wg.Done()
			}
		}(i)

	}

	wg.Wait()

	totalT = time.Now().UnixNano() - totalT
	fmt.Printf("took %f ms for %d requests\n", float64(totalT)/1000000, n*m)

	totalD := make([]int64, 0, n*m)
	for _, k := range d {
		totalD = append(totalD, k...)
	}
	totalD2 := make([]float64, 0, n*m)
	for _, k := range totalD {
		totalD2 = append(totalD2, float64(k))
	}

	mean, _ := stats.Mean(totalD2)
	median, _ := stats.Median(totalD2)
	max, _ := stats.Max(totalD2)
	min, _ := stats.Min(totalD2)
	p99, _ := stats.Percentile(totalD2, 99.9)

	fmt.Printf("sent     requests    : %d\n", n*m)
	fmt.Printf("received requests    : %d\n", atomic.LoadUint64(&trans))
	fmt.Printf("received requests_OK : %d\n", atomic.LoadUint64(&transOK))
	fmt.Printf("throughput  (TPS)    : %d\n", int64(n*m)*1000000000/totalT)
	fmt.Printf("mean: %.f ns, median: %.f ns, max: %.f ns, min: %.f ns, p99.9: %.f ns\n", mean, median, max, min, p99)
	fmt.Printf("mean: %d ms, median: %d ms, max: %d ms, min: %d ms, p99: %d ms\n", int64(mean/1000000), int64(median/1000000), int64(max/1000000), int64(min/1000000), int64(p99/1000000))

	common.InitSignal(survivalTimeout)
}
