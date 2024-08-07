package grpcpool

import (
	"errors"
)

var (
	ErrInvalidPoolSetting = errors.New("grpc pool: client pool setting invalid")
	ErrInvalidFactory     = errors.New("grpc pool: invalid generate factory")
	ErrClosed             = errors.New("grpc pool: client pool is closed")
	ErrConnectionLess     = errors.New("grpc pool: no available connection")
	ErrConnectionFull     = errors.New("grpc pool: connection pool is full")
)

type (
	Pool interface {
		Get() (*Conn, error)

		Close()

		Len() int
	}
)
