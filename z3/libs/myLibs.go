package libs

import (
	"sync"
)

type Semaphore struct {
	channel  chan struct{}
	capacity int
}

func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		channel:  make(chan struct{}, capacity),
		capacity: capacity}
}

func (s *Semaphore) Acquire() {
	s.channel <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.channel
}

type Monitor struct {
	mutex sync.Mutex
	cond  *sync.Cond
}

func NewMonitor() *Monitor {
	m := &Monitor{}
	m.cond = sync.NewCond(&m.mutex)
	return m
}

func (m *Monitor) Enter() {
	m.mutex.Lock()
}

func (m *Monitor) Exit() {
	m.mutex.Unlock()
}

func (m *Monitor) Wait() {
	m.cond.Wait()
}

func (m *Monitor) Signal() {
	m.cond.Signal()
}

func (m *Monitor) Broadcast() {
	m.cond.Broadcast()
}
