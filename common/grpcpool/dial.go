package grpcpool

import (
	"context"
	"time"

	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc/backoff"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

const (
	dialTimeout = 5 * time.Second
	// backoffMaxDelay provided maximum delay when backing off after failed connection attempts.
	backoffMaxDelay = 3 * time.Second
	maxSendMsgSize  = 4 << 20
	maxRecvMsgSize  = 4 << 20
	// keepAliveTimeout is the duration of time for which the client waits after having
	// pinged for keepalive check and if no activity is seen even after that the connection
	// is closed.
	keepAliveTimeout = time.Duration(3) * time.Second
	// keepAliveTime is the duration of time after which if the client doesn't see
	// any activity it pings the server to see if the transport is still alive.
	keepAliveTime = time.Duration(10) * time.Second
)

type (
	Factory func(addr string) (*grpc.ClientConn, error)

	Options struct {
		dial Factory

		initCount int
		maxCount  int
	}
)

var defaultOptions = Options{
	dial:      Dial,
	initCount: 5,
	maxCount:  30,
}

func Dial(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	return grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.DefaultConfig}),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(maxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxRecvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepAliveTime,
			Timeout:             keepAliveTimeout,
			PermitWithoutStream: true,
		}))

}

type OptFn = func(option *Options)

func WithDial(dial Factory) OptFn {
	return func(option *Options) {
		option.dial = dial
	}
}

func WithInitCount(c int) OptFn {
	return func(option *Options) {
		option.initCount = c
	}
}

func WithMaxCount(c int) OptFn {
	return func(option *Options) {
		option.maxCount = c
	}
}
