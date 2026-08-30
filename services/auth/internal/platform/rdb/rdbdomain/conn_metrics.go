package rdbdomain

import "sync"

type ConnMetric struct {
	mu         sync.RWMutex
	activeConn int64
}

func newConnMetric() *ConnMetric {
	return &ConnMetric{
		mu:         sync.RWMutex{},
		activeConn: 0,
	}
}

func (c *ConnMetric) onCreate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.activeConn++
}

func (c *ConnMetric) getActiveConn() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.activeConn
}
