package postgres

import (
	"sync"
)

type Metrics struct {
	onCreate    int64
	onAcquire   int64
	onRelease   int64
	onDestroyed int64
	mu          *sync.RWMutex
}

func (m *Metrics) OnCreate() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onCreate += 1
}

func (m *Metrics) OnAcquire() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onAcquire += 1
}

func (m *Metrics) OnRelease() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onRelease += 1
}

func (m *Metrics) OnDestroyed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onDestroyed += 1
}

func (m *Metrics) GetOnCreate() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.onCreate
}

func (m *Metrics) GetOnAcquire() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.onAcquire
}

func (m *Metrics) GetOnRelease() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.onRelease
}

func (m *Metrics) GetOnDestroyed() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.onDestroyed
}
