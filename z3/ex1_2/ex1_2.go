package main

import (
	"fmt"
	"sync"
	"time"
	"z3/libs"
)

const numPhilosophers = 5

type Message struct {
	idx   int
	state int
}

type Philosopher struct {
	id        int
	leftFork  *libs.Semaphore
	rightFork *libs.Semaphore
}

func newPhilosopher(id int, leftFork *libs.Semaphore, rightFork *libs.Semaphore) *Philosopher {
	return &Philosopher{
		id,
		leftFork,
		rightFork,
	}
}

func (p *Philosopher) dine(messages chan Message, wg *sync.WaitGroup, diner *chan int) {
	defer wg.Done()
	for {
		*diner <- p.id
		p.rightFork.Acquire()
		p.leftFork.Acquire()
		messages <- Message{p.id, 0}
		time.Sleep(1 * time.Second)
		p.rightFork.Release()
		p.leftFork.Release()
		<-*diner
		messages <- Message{p.id, 1}
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
	forks := make([]*libs.Semaphore, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = libs.NewSemaphore(1)
	}
	philosophers := make([]*Philosopher, numPhilosophers)
	messages := make(chan Message)
	go printMessages(messages)
	diner := make(chan int, numPhilosophers-1)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = newPhilosopher(i, forks[i], forks[(i+1)%numPhilosophers])
		go philosophers[i].dine(messages, &wg, &diner)
	}
	wg.Wait()
}
