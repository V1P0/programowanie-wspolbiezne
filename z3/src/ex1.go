package main

import (
	"fmt"
	"sync"
	"time"
	"z3/libs"
)

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

func (p *Philosopher) dine(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		p.rightFork.Acquire()
		p.leftFork.Acquire()
		fmt.Printf("Philosopher %d is eating\n", p.id)
		time.Sleep(1 * time.Second)
		p.rightFork.Release()
		p.leftFork.Release()
		fmt.Printf("Philosopher %d is thinking\n", p.id)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	const numPhilosophers = 5
	var wg sync.WaitGroup
	wg.Add(numPhilosophers)
	forks := make([]*libs.Semaphore, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = libs.NewSemaphore(1)
	}
	philosophers := make([]*Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = newPhilosopher(i, forks[i], forks[(i+1)%numPhilosophers])
		go philosophers[i].dine(&wg)
	}
	wg.Wait()
}
