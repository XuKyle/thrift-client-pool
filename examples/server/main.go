package main

import (
	"github.com/XuKyle/thrift-client-pool/api"
	"github.com/XuKyle/thrift-client-pool/pkg/handler"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
)

func main() {
	addr := "127.0.0.1:9999"
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())

	var transport thrift.TServerTransport
	var err error

	transport, err = thrift.NewTServerSocket(addr)
	if err != nil {
		log.Println(err)
	}

	handler := &handler.AddHandler{}
	addServiceProcessor := api.NewAddServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(addServiceProcessor, transport, transportFactory, protocolFactory)
	server.Serve()
}
