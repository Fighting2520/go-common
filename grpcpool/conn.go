package grpcpool

import (
	"sync"

	"google.golang.org/grpc"
)

type (
	Conn struct {
		*grpc.ClientConn
		mu      sync.RWMutex
		pool    *chanPool
		useless bool
	}
)

func (c *Conn) isUseless() bool {
	return c.useless
}

func (c *Conn) MarkUseless() {
	c.mu.Lock()
	c.useless = true
	c.mu.Unlock()
}

func newConn(c *chanPool) (*Conn, error) {
	conn, err := c.factory(c.addr)
	if err != nil {
		return nil, err
	}
	return &Conn{ClientConn: conn, pool: c}, nil
}

func (c *Conn) reset() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ClientConn != nil {
		c.ClientConn.Close()
	}
	conn, err := c.pool.factory(c.pool.addr)
	if err != nil {
		return err
	}
	c.useless = false
	c.ClientConn = conn
	return nil
}

func (c *Conn) Close() error {
	if c == nil {
		return nil
	}
	if c.isUseless() {
		return c.ClientConn.Close()
	}
	// valid connection
	return c.pool.put(c)
}
