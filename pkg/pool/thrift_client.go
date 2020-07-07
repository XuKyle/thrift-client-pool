package pool

import (
	"github.com/apache/thrift/lib/go/thrift"
	"time"
)

// Thrift客户端
type IdleClient struct {
	// Thrift传输层，封装了底层连接建立、维护、关闭、数据读写等细节
	Transport thrift.TTransport
	// 真正的Thrift客户端，业务创建传入
	RawClient interface{}
}

// 检测连接是否有效
func (client *IdleClient) Check() bool {
	if client.Transport == nil || client.RawClient == nil {
		return false
	}
	return client.Transport.IsOpen()
}

// 封装了Thrift客户端
type idleConn struct {
	// 空闲Thrift客户端
	c *IdleClient
	// 最近一次放入空闲队列的时间
	t time.Time
}

// Thrift客户端创建方法，留给业务去实现
type ThriftDial func(addr []string, connTimeout time.Duration) (*IdleClient, error)

//关闭Thrift客户端，留给业务实现
type ThriftClientClose func(client *IdleClient) error
