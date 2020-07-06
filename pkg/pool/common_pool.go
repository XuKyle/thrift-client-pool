package pool

import "sync"

type Pool interface {
	initConn()
	createIdleConn(wg *sync.WaitGroup)

	Get() (*IdleClient, error)
	Put(client *IdleClient) error
	Reconnect(client *IdleClient) (newClient *IdleClient, err error)
	CloseConn(client *IdleClient)
	ClearConn()
	Release()
	Recover()

	CheckTimeout()
	GetIdleCount() uint32
	GetConnCount() int32
}

type Agent interface {
	Init(pool *ThriftPool)

	Do(do func(rawClient interface{}) error) error

	getClient() (*IdleClient, error)
	releaseClient(client *IdleClient) error
	reconnect(client *IdleClient) (newClient *IdleClient, err error)
	closeClient(client *IdleClient)

	Release()
	GetIdleCount() uint32
	GetConnCount() int32
}
