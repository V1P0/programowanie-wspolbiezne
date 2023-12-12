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

func newMonitor() *Monitor {
	m := &Monitor{}
	m.cond = sync.NewCond(&m.mutex)
	return m
}

func (m *Monitor) enter() {
	m.mutex.Lock()
}

func (m *Monitor) exit() {
	m.mutex.Unlock()
}

func (m *Monitor) wait() {
	m.cond.Wait()
}

func (m *Monitor) signal() {
	m.cond.Signal()
}

func (m *Monitor) broadcast() {
	m.cond.Broadcast()
}
