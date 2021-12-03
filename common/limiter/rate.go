package limiter

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// Limit 定义某些事件的最大频率。
// Limit 表示为每秒的事件数。
// zero Limit 不允许任何事件
type Limit float64

// Inf 无限速 Limit 允许任何事件（即使burst是0）
const Inf = Limit(math.MaxFloat64)

// Every 如果interval是limit的时间间隔, 每秒的令牌数量
// Every interval 每个Limit 的最小时间间隔
func Every(interval time.Duration) Limit {
	if interval <= 0 {
		return Inf
	}
	return 1 / Limit(interval.Seconds())
}

type Limiter struct {
	mu     sync.Mutex
	limit  Limit
	burst  int
	tokens float64
	// last 指的是 tokens 最近一次被更新的时间
	last      time.Time
	lastEvent time.Time
}

// Limit 返回最大的总体速率限制
func (lim *Limiter) Limit() Limit {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.limit
}

// Burst 返回最大的突发大小, 即 burst 字段
func (lim *Limiter) Burst() int {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.burst
}

// NewLimiter 实例化一个限制器
func NewLimiter(r Limit, b int) *Limiter {
	return &Limiter{
		limit: r,
		burst: b,
	}
}

func (lim *Limiter) Allow() bool {
	return lim.AllowN(time.Now(), 1)
}

// AllowN 判断指定时间是否能够允许N个时间
func (lim *Limiter) AllowN(now time.Time, n int) bool {
	return lim.reserveN(now, n, 0).ok
}

// Reservation 保存 Limiter 允许在延迟后发生的事件的信息。
// A Reservation 可以被取消， 可以使 Limiter 允许更多事件
type Reservation struct {
	ok        bool
	lim       *Limiter
	tokens    int
	timeToAct time.Time
	limit     Limit
}

func (r *Reservation) Ok() bool {
	return r.ok
}

func (r *Reservation) Delay() time.Duration {
	return r.DelayFrom(time.Now())
}

const InfDuration = time.Duration(1<<63 - 1)

func (r *Reservation) DelayFrom(now time.Time) time.Duration {
	if !r.ok {
		return InfDuration
	}
	delay := r.timeToAct.Sub(now)
	if delay < 0 {
		return 0
	}
	return delay
}

func (r *Reservation) Cancel() {
	r.CancelAt(time.Now())
}

func (r *Reservation) CancelAt(now time.Time) {
	if !r.ok {
		return
	}
	r.lim.mu.Lock()
	defer r.lim.mu.Unlock()

	if r.lim.limit == Inf || r.tokens == 0 || r.timeToAct.Before(now) {
		return
	}

	// 计算要恢复的令牌
	// The duration between lim.lastEvent and r.timeToAct tells us how many tokens were reserved after r was obtained. These tokens should not be restored.
	restoreTokens := float64(r.tokens) - r.limit.tokensFromDuration(r.lim.lastEvent.Sub(r.timeToAct))
	if restoreTokens < 0 {
		return
	}
	var tokens float64
	now, _, tokens = r.lim.advance(now)
	// 计算新的令牌数量
	tokens += restoreTokens
	if burst := float64(r.lim.burst); tokens > burst {
		tokens = burst // 相当于一个令牌桶，桶中最多能放 burst 个令牌
	}
	// 更新状态
	r.lim.last = now
	r.lim.tokens = tokens
	if r.timeToAct == r.lim.lastEvent {
		prevEvent := r.timeToAct.Add(r.limit.durationFromTokens(float64(-r.tokens)))
		if !prevEvent.Before(now) {
			r.lim.lastEvent = prevEvent
		}
	}
}

func (lim *Limiter) Reserve() *Reservation {
	return lim.ReserveN(time.Now(), 1)
}

func (lim *Limiter) ReserveN(now time.Time, n int) *Reservation {
	r := lim.reserveN(now, n, InfDuration)
	return &r
}

func (lim *Limiter) Wait(ctx context.Context) error {
	return lim.WaitN(ctx, 1)
}

// WaitN 阻塞，直到 Limiter 允许n个事件发生
func (lim *Limiter) WaitN(ctx context.Context, n int) error {
	lim.mu.Lock()
	limit := lim.limit
	burst := lim.burst
	lim.mu.Unlock()

	if n > burst && limit != Inf {
		return fmt.Errorf("rate: Wait(n=%d) exceeds limiter's burst %d", n, burst)
	}
	// Check if ctx is already cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	now := time.Now()
	waitLimit := InfDuration // 等待时间的阈值
	if deadline, ok := ctx.Deadline(); ok {
		waitLimit = deadline.Sub(now) // 如果设置了超时时间，则获取超时时间作为等待时间的阈值
	}

	r := lim.reserveN(now, n, waitLimit)
	if !r.Ok() {
		return fmt.Errorf("rate: Wait(n=%d) would exceed context deadline", n)
	}
	delay := r.DelayFrom(now)
	if delay == 0 {
		return nil
	}
	t := time.NewTimer(delay)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		r.Cancel()
		return ctx.Err()
	}
}

func (lim *Limiter) SetLimit(newLimit Limit) {
	lim.SetLimitAt(time.Now(), newLimit)
}

func (lim *Limiter) SetLimitAt(now time.Time, newLimit Limit) {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	var tokens float64
	now, _, tokens = lim.advance(now)

	lim.last = now
	lim.tokens = tokens
	lim.limit = newLimit
}

func (lim *Limiter) SetBurst(newBurst int) {
	lim.SetBurstAt(time.Now(), newBurst)
}

// SetBurstAt 重新设置一个新突发大小
func (lim *Limiter) SetBurstAt(now time.Time, newBurst int) {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	var tokens float64
	now, _, tokens = lim.advance(now)

	lim.last = now
	lim.tokens = tokens
	lim.burst = newBurst
}

// maxFutureReserve 最大预约等待时间
func (lim *Limiter) reserveN(now time.Time, n int, maxFutureReserve time.Duration) Reservation {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	if lim.limit == Inf {
		return Reservation{
			ok:        true,
			lim:       lim,
			tokens:    n,
			timeToAct: now,
		}
	} else if lim.limit == 0 {
		// zero Limiter
		var ok bool
		if lim.burst > n {
			ok = true
			lim.burst -= n
		}
		return Reservation{
			ok:        ok,
			lim:       lim,
			tokens:    lim.burst,
			timeToAct: now,
		}
	}
	now, last, tokens := lim.advance(now)
	// 计算请求产生的剩余令牌数量。
	tokens -= float64(n)
	// 计算等待时间
	var waitDuration time.Duration
	if tokens < 0 {
		waitDuration = lim.limit.durationFromTokens(-tokens)
	}
	ok := n <= lim.burst && waitDuration <= maxFutureReserve
	r := Reservation{
		ok:    ok,
		lim:   lim,
		limit: lim.limit,
	}
	if ok {
		r.tokens = n
		r.timeToAct = now.Add(waitDuration)

		lim.last = now
		lim.tokens = tokens
		lim.lastEvent = r.timeToAct
	} else {
		lim.last = last
	}
	return r
}

func (lim *Limiter) advance(now time.Time) (newNow, newLast time.Time, newTokens float64) {
	last := lim.last
	if now.Before(last) {
		last = now
	}
	// 计算由于经过的时间而产生的新的令牌数量
	elapsed := now.Sub(last)
	delta := lim.limit.tokensFromDuration(elapsed)
	tokens := lim.tokens + delta
	if burst := float64(lim.burst); tokens > burst {
		tokens = burst // 相当于一个令牌桶，桶中最多能放 burst 个令牌
	}
	return now, last, tokens
}

// 累积到tokens数量的令牌需要的时间
func (limit Limit) durationFromTokens(tokens float64) time.Duration {
	if limit < 0 {
		return InfDuration
	}
	seconds := tokens / float64(limit)
	return time.Duration(seconds * float64(time.Second))
}

// tokensFromDuration 从 持续时间 d 时间内能够积累的令牌数
func (limit Limit) tokensFromDuration(d time.Duration) float64 {
	if limit <= 0 {
		return 0
	}

	return d.Seconds() * float64(limit)
}
