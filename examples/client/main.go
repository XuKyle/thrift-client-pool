package main

import (
	"context"
	"fmt"
	"github.com/XuKyle/thrift-client-pool/api"
	"github.com/apache/thrift/lib/go/thrift"
	"math/rand"
)

func main() {
	var transport thrift.TTransport
	transport, _ = thrift.NewTSocket(":9999")
	protocolTransport := thrift.NewTBinaryProtocolTransport(thrift.NewTFramedTransport(transport))
	tClient := thrift.NewTStandardClient(protocolTransport, protocolTransport)

	client := api.NewAddServiceClient(tClient)

	if err := transport.Open(); err != nil {
		fmt.Println("error open transport", err)
		return
	}
	defer transport.Close()

	num1 := rand.Int31n(1000)
	num2 := rand.Int31n(1000)

	result, err := client.Add(context.Background(), num1, num2)
	if err != nil {
		fmt.Println("err : add", err)
	}

	fmt.Println(num1, "+", num2, "=", result)
}
