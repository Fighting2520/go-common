package limiter

import "testing"

func TestTokenBucketLimiter_Allow(t *testing.T) {
	runNodeLimit(t, NewNodeTokenLimiter(10, 1), []allow{
		{t0, 1, true},
		{t0, 1, false},
		{t0, 1, false},
		{t1, 1, true},
		{t1, 1, false},
		{t1, 1, false},
		{t2, 2, false}, // burst size is 1, so n=2 always fails
		{t2, 1, true},
		{t2, 1, false},
		{t2, -10, true}, // always return true when n < 0
		{t2, 1, false},
	})
}

func runNodeLimit(t *testing.T, lim *NodeTokenLimiter, allows []allow) {
	t.Helper()
	for i, allow := range allows {
		ok := lim.AllowN(allow.t, allow.n)
		if ok != allow.ok {
			t.Errorf("step %d: lim.AllowN(%v, %v) = %v want %v",
				i, allow.t, allow.n, ok, allow.ok)
		}
	}
}
