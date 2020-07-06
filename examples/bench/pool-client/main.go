package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/XuKyle/thrift-client-pool/api"
	"github.com/XuKyle/thrift-client-pool/pkg/pool"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	addr            = "127.0.0.1:9999"
	connTimeout     = time.Second * 2
	idleTimeout     = time.Second * 120
	timeout         = time.Second * 10
	maxConn         = int32(100)
	bct             = context.Background()
	delay           int64 // 单位微妙
	successCount    = int64(0)
	failCount       = int64(0)
	thriftPoolAgent *pool.ThriftPoolAgent
)

// 初始化Thrift连接池代理
func init() {
	config := &pool.ThriftPoolConfig{
		Addr:        addr,
		MaxConn:     maxConn,
		ConnTimeout: connTimeout,
		IdleTimeout: idleTimeout,
		Timeout:     timeout,
	}
	thriftPool := pool.NewThriftPool(config, thriftDial, thriftClose)
	thriftPoolAgent = pool.NewThriftPoolAgent()
	thriftPoolAgent.Init(thriftPool)
}

func thriftDial(addr string, duration time.Duration) (*pool.IdleClient, error) {
	var transport thrift.TTransport
	transport, err := thrift.NewTSocketTimeout(addr, connTimeout)
	if err != nil {
		return nil, err
	}

	protocolTransport := thrift.NewTBinaryProtocolTransport(thrift.NewTFramedTransport(transport))
	tClient := thrift.NewTStandardClient(protocolTransport, protocolTransport)
	client := api.NewAddServiceClient(tClient)

	if err := transport.Open(); err != nil {
		fmt.Println("error open transport", err)
		return nil, err
	}

	return &pool.IdleClient{
		Transport: transport,
		RawClient: client,
	}, nil
}

func thriftClose(client *pool.IdleClient) error {
	if client == nil {
		return nil
	}
	return client.Transport.Close()
}

func startPprof() {
	http.ListenAndServe("0.0.0.0:9998", nil)
}

func success(st time.Time) {
	atomic.AddInt64(&delay, time.Now().Sub(st).Microseconds())
	atomic.AddInt64(&successCount, 1)
}

func fail() {
	atomic.AddInt64(&failCount, 1)
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

func dialApi() (result int32, err error) {
	st := time.Now()
	err = thriftPoolAgent.Do(func(rawClient interface{}) error {
		client, ok := rawClient.(*api.AddServiceClient)
		if !ok {
			return errors.New("unknown client type")
		}

		var err2 error
		num1 := rand.Int31n(1000)
		num2 := rand.Int31n(1000)

		result, err2 = client.Add(bct, num1, num2)
		return err2
	})

	if err != nil {
		fail()
	} else {
		success(st)
	}

	return
}

func main() {
	go startPprof()

	//start()
	for i := 0; i < 100; i++ {
		go start()
	}
	time.Sleep(time.Second * 10)
	avgQps := float64(successCount) / float64(10)
	avgDelay := float64(delay) / float64(successCount) / 1000
	log.Println(fmt.Sprintf("总运行时间：600s, 并发协程数：100，平均吞吐量：%v，平均延迟（ms）：%v，总成功数：%d，总失败数：%d",
		avgQps, avgDelay, successCount, failCount))
}
