package limiter

import (
	"golang.org/x/time/rate"
	"time"
)

type (
	NodeTokenLimiter struct {
		limit float64
		burst int
		lim   *rate.Limiter
	}
)

func NewNodeTokenLimiter(limit float64, burst int) *NodeTokenLimiter {
	return &NodeTokenLimiter{
		limit: limit,
		burst: burst,
		lim:   rate.NewLimiter(rate.Limit(limit), burst),
	}
}

func (nl *NodeTokenLimiter) Allow() bool {
	return nl.AllowN(time.Now(), 1)
}

// AllowN 重写 rate.AllowN方法 修复人为可通过设置n<0的数来反向增加令牌
func (nl *NodeTokenLimiter) AllowN(now time.Time, n int) bool {
	if n < 0 {
		return true
	}
	return nl.lim.AllowN(now, n)
}
