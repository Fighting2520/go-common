package grpcpool

import (
	"context"
	"sync"
	"time"
)

type (
	chanPool struct {
		addr  string
		conns chan *Conn
		mu    sync.RWMutex

		factory Factory
	}
)

func NewPool(addr string, opts ...OptFn) (Pool, error) {
	op := defaultOptions
	for _, fn := range opts {
		fn(&op)
	}
	if op.initCount < 0 || op.maxCount <= 0 || op.initCount > op.maxCount {
		return nil, ErrInvalidPoolSetting
	}
	c := &chanPool{
		addr:    addr,
		conns:   make(chan *Conn, op.maxCount),
		factory: op.dial,
	}

	for i := 0; i < op.initCount; i++ {
		conn, err := op.dial(addr)
		if err != nil {
			c.Close()
			return nil, ErrInvalidFactory
		}
		c.conns <- &Conn{ClientConn: conn}
	}

	return c, nil
}

func (c *chanPool) Get() (*Conn, error) {
	if c.conns == nil {
		return nil, ErrClosed
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(500)*time.Millisecond)
	defer cancelFunc()
	select {
	case conn := <-c.conns:
		if conn == nil {
			// 无效连接
			return c.new()
		}
		if conn.isUseless() {
			err := conn.reset()
			if err != nil {
				return nil, err
			}
		}
		return conn, nil
	case <-ctx.Done():
		return nil, ErrConnectionLess
	}
}

func (c *chanPool) new() (*Conn, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return newConn(c)
}

func (c *chanPool) reset() {
	c.mu.Lock()

}

func (c *chanPool) Close() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		if conn != nil {
			conn.Close()
		}
	}
}

func (c *chanPool) Len() int {
	return len(c.conns)
}

func (c *chanPool) put(conn *Conn) error {
	c.mu.Lock()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(500)*time.Millisecond)
	defer cancelFunc()
	select {
	case c.conns <- conn:
		// valid op
	case <-ctx.Done():
		return ErrConnectionFull
	}
	return nil
}
