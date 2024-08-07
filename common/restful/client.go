package restful

import (
	"net/http"
	"time"
)

const (
	defaultTimeout = 5 * time.Second
)

type Client struct {
	option
	client *http.Client
}

func NewClient(fns ...OptFn) *Client {
	var opt option
	for _, fn := range fns {
		fn(&opt)
	}

	if opt.timeout <= 0 {
		opt.timeout = defaultTimeout
	}

	client := http.DefaultClient
	client.Timeout = opt.timeout
	return &Client{
		option: opt,
		client: client,
	}
}
