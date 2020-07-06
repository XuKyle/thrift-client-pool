package pool

import "time"

// 连接池配置
type ThriftPoolConfig struct {
	// Thrfit Server端地址
	Addr string
	// 最大连接数
	MaxConn int32
	// 创建连接超时时间
	ConnTimeout time.Duration
	// 空闲客户端超时时间，超时主动释放连接，关闭客户端
	IdleTimeout time.Duration
	// 获取Thrift客户端超时时间
	Timeout time.Duration
	// 获取Thrift客户端失败重试间隔
	interval time.Duration
}
