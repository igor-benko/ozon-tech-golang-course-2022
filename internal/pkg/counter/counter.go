package counter

import "sync/atomic"

type Counter struct {
	count uint64
}

func (c *Counter) Inc() {
	atomic.AddUint64(&c.count, 1)
}

func (c *Counter) Get() uint64 {
	return c.count
}
