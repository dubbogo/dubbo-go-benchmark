package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

import (
	"github.com/montanaflynn/stats"
)

import (
	"github.com/dubbo/go-for-apache-dubbo/config"
	_ "github.com/dubbo/go-for-apache-dubbo/protocol/jsonrpc"
	_ "github.com/dubbo/go-for-apache-dubbo/registry/protocol"

	_ "github.com/dubbo/go-for-apache-dubbo/filter/impl"

	_ "github.com/dubbo/go-for-apache-dubbo/cluster"
	_ "github.com/dubbo/go-for-apache-dubbo/cluster/loadbalance"
	_ "github.com/dubbo/go-for-apache-dubbo/common/proxy/proxy_factory"

	_ "github.com/dubbo/go-for-apache-dubbo/cluster/cluster_impl"

	_ "github.com/dubbo/go-for-apache-dubbo/registry/zookeeper"
)

// they are necessary:
// 		export CONF_CONSUMER_FILE_PATH="xxx"
// 		export APP_LOG_CONF_FILE="xxx"
var concurrency = flag.Int("c", 1, "concurrency")
var total = flag.Int("n", 1, "total requests for all clients")

// package 4096(req, more) 4141(rsp)
const ARG = "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"

// package 1024(req) 1069(rsp)
const ARG1 = "UjSNFgAIBgDCggDyHy3fB5BwT1COQKDHRkB6xT+8BDeOUcctcqPdiO/3NNNtxVDgt6E+i/zb7OMcLADgAe5UjSNFCgAIBgDyHy6hN5GkZndiWJ48ERFRFUKqq9iN4/zb7OMc8386x2Al23Z5a3fB5BwT1C4OQKDHRkB6xT+8BDeOUcctcqPdiO/3NNNtxVDgt6E+i/zb7OMcLADgAe5UjSNFCgAIBgDyHy6hN5GkZndiWJ48s9OMVORckH/AABBNuP8fp73U9U8uskqaVta3+wtyPz/ALD3aP4q9eHSbrWsG+oJHH+HA/1ufaj068TTHy6DncThaWdrkaAjNlsbjV6V0qR/gLj+vtYAWBJHSQNpoNVM9MOcnqalXhRBTr4+dDiEViW1B6KLhEBJpnpReEuNViQVAJAjNlsbjV6V0qR/gLj+vtYAWBJHSQNpoNVM9MOcnqa+lXhRBTr4+XUgKLXOkfQkkAg/4Gw9P+8e/Ak9J2SmFB6TOczdDi4JaipmjREViWawSwF78KR/tr+9Enp6O3egJWg6MN8ffjPl+0J6HfPZNBNi9iN46vD7Sqo9J2SmFB6TOczdDi4JaipmjREViWawSwF78KR/tr+9Enp6O3egJWg6MN8ffjPl+0J6HfPZNBNi9iN"

// package 600(req) 645(rsp)
const ARG2 = "8BDeOUcct2cqPdiO3NNNtxVDgt6E+i/zb7OMcLADgAe5UjSNFCgAIBgDyHy6hN5GkZndiWJ48T/s9OMVORckH/AABBNuP8fp73U9U8uskqaVta3+wtyPz/ALD3aP4q9eHSbrWsG+oJHH+HA/1ufaj068TTHy6DncThaWdrkaAjNlsbjV6V0qR/gLj+vtYAWBJHSQNpoNVM9MOcnqalXhRBTr4+dDiEViW1"

func main() {
	flag.Parse()

	conc, tn, err := checkArgs(*concurrency, *total)
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

	conMap, _ := config.Load()
	if conMap == nil {
		panic("conMap is nil")
	}

	time.Sleep(3e9)

	var startWg sync.WaitGroup
	startWg.Add(n)

	var trans uint64
	var transOK uint64

	d := make([][]int64, n, n)

	//it contains warmup time but we can ignore it
	totalT := time.Now().UnixNano()
	for i := 0; i < n; i++ {
		dt := make([]int64, 0, m)
		d = append(d, dt)

		go func(i int) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered in f", r)
				}
			}()

			//warmup
			//for j := 0; j < 5; j++ {
			//	user := &JsonRPCUser{}
			//	err := conMap["com.ikurento.user.UserProvider"].GetRPCService().(*UserProvider).GetUser(context.TODO(), []interface{}{"A003"}, user)
			//	if err != nil {
			//		fmt.Println(err)
			//	}
			//}

			startWg.Done()
			startWg.Wait()

			for j := 0; j < m; j++ {
				t := time.Now().UnixNano()
				user := &User{}
				err := conMap["com.ikurento.user.UserProvider"].GetRPCService().(*UserProvider).GetUser(context.TODO(), []interface{}{"A003", ARG}, user)

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
}

// checkArgs check concurrency and total request count.
func checkArgs(c, n int) (int, int, error) {
	if c < 1 {
		fmt.Printf("c < 1 and reset c = 1")
		c = 1
	}
	if n < 1 {
		fmt.Printf("n < 1 and reset n = 1")
		n = 1
	}
	if c > n {
		return c, n, errors.New("c must be set <= n")
	}
	return c, n, nil
}
