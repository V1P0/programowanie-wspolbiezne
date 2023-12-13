package main

import (
	"fmt"
	"sync"
	"time"
)

const numPhilosophers = 5

type Message struct {
	idx   int
	state int
}

type Monitor struct {
	forks   [numPhilosophers]int
	okToEat [numPhilosophers]chan int
	lock    chan int
}

func (m *Monitor) takeFork(i int, messages chan Message) {
	m.lock <- 1
	if m.forks[i] != 2 {
		<-m.lock
		for {
			<-m.okToEat[i]
			m.lock <- 1
			if m.forks[i] == 2 {
				break
			}
			<-m.lock
		}
	}

	m.forks[(i+1)%numPhilosophers] -= 1
	m.forks[(i-1+numPhilosophers)%numPhilosophers] -= 1
	messages <- Message{i, 0}
	<-m.lock
}

func (m *Monitor) releaseFork(i int, messages chan Message) {
	m.lock <- 1
	m.forks[(i+1)%numPhilosophers] += 1
	m.forks[(i-1+numPhilosophers)%numPhilosophers] += 1

	if m.forks[(i+1)%numPhilosophers] == 2 {
		select {
		case m.okToEat[(i+1)%numPhilosophers] <- 1:
		default:
		}
	}

	if m.forks[(i-1+numPhilosophers)%numPhilosophers] == 2 {
		select {
		case m.okToEat[(i-1+numPhilosophers)%numPhilosophers] <- 1:
		default:
		}
	}
	messages <- Message{i, 1}

	<-m.lock
}

type Philosopher struct {
	id      int
	monitor *Monitor
}

func newPhilosopher(id int, monitor *Monitor) *Philosopher {
	return &Philosopher{
		id,
		monitor,
	}
}

func (p *Philosopher) dine(messages chan Message, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		p.monitor.takeFork(p.id, messages)
		time.Sleep(1 * time.Second)
		p.monitor.releaseFork(p.id, messages)
		time.Sleep(1 * time.Second)
	}
}

func printMessages(messages chan Message) {
	var current_list [numPhilosophers]int

	for {
		message := <-messages
		current_list[message.idx] = message.state

		fmt.Printf("List:\n")
		for i := 0; i < numPhilosophers; i++ {
			if current_list[i] == 1 {
				fmt.Printf("(w_%d, f_%d, w_%d)\n", i, i, (i+1)%numPhilosophers)
			}
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(numPhilosophers)
	philosophers := make([]*Philosopher, numPhilosophers)
	messages := make(chan Message)
	var monitor Monitor
	for i := 0; i < numPhilosophers; i++ {
		monitor.okToEat[i] = make(chan int, 1)
	}

	for i := 0; i < numPhilosophers; i++ {
		monitor.forks[i] = 2
	}

	monitor.lock = make(chan int, 1)
	go printMessages(messages)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = newPhilosopher(i, &monitor)
		go philosophers[i].dine(messages, &wg)
	}
	wg.Wait()
}
