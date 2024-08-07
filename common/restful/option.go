package restful

import "time"

type (
	option struct {
		username, password string
		timeout            time.Duration
	}

	OptFn = func(option *option)
)

func WithTimeout(timeout time.Duration) OptFn {
	return func(option *option) {
		option.timeout = timeout
	}
}

func WithAuth(username, password string) OptFn {
	return func(option *option) {
		option.username = username
		option.password = password
	}
}
