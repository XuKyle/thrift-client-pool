package main

import (
	"context"
	"fmt"
	"github.com/XuKyle/thrift-client-pool/api"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync/atomic"
	"time"
)

var (
	addr         = "127.0.0.1:9999"
	connTimeout  = time.Second * 2
	bct          = context.Background()
	delay        int64 // 单位微妙
	successCount = int64(0)
	failCount    = int64(0)
)

func main() {
	go startPprof()
	for i := 0; i < 100; i++ {
		go start()
	}
	time.Sleep(time.Second * 60)
	avgQps := float64(successCount) / float64(60)
	avgDelay := float64(delay) / float64(successCount) / 1000
	log.Println(fmt.Sprintf("总运行时间：600s, 并发协程数：100，平均吞吐量：%v，平均延迟（ms）：%v，总成功数：%d，总失败数：%d",
		avgQps, avgDelay, successCount, failCount))
}

func startPprof() {
	http.ListenAndServe("0.0.0.0:9998", nil)
}

func start() {
	for {
		if resp, err := dialApi(); err != nil {
			log.Println(err)
		} else {
			log.Println(resp)
		}
	}
}

func dialApi() (int32, error) {
	st := time.Now()

	var transport thrift.TTransport
	transport, err := thrift.NewTSocketTimeout(addr, connTimeout)
	if err != nil {
		fail()
		return 0, err
	}

	if err := transport.Open(); err != nil {
		fail()
		return 0, err
	}
	defer transport.Close()

	protocolTransport := thrift.NewTBinaryProtocolTransport(thrift.NewTFramedTransport(transport))
	tClient := thrift.NewTStandardClient(protocolTransport, protocolTransport)

	client := api.NewAddServiceClient(tClient)

	num1 := rand.Int31n(1000)
	num2 := rand.Int31n(1000)

	resp, err := client.Add(bct, num1, num2)
	if err != nil {
		fail()
		return 0, nil
	} else {
		success(st)
		return resp, nil
	}
}

func success(st time.Time) {
	atomic.AddInt64(&delay, time.Now().Sub(st).Microseconds())
	atomic.AddInt64(&successCount, 1)
}

func fail() {
	atomic.AddInt64(&failCount, 1)
}
