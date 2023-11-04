package pool

import (
	"fmt"
	"sync"
)

type Worker struct {
	readyWg *sync.WaitGroup

	tasksChannel chan func()
	closeChannel chan bool
}

func (w *Worker) Run() {
	w.readyWg.Done()

	for {
		select {
		case <-w.closeChannel:
			close(w.tasksChannel)
			return
		case task := <-w.tasksChannel:
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err) // change to log.Error
				}
			}()
			task()
		}
	}
}

func NewWorker(wg *sync.WaitGroup) *Worker {
	return &Worker{
		readyWg:      wg,
		tasksChannel: make(chan func()),
		closeChannel: make(chan bool),
	}
}

type Pool struct {
	workers          []*Worker
	latestUsedWorker int
	submitLock       sync.Mutex
}

func (p *Pool) Submit(f func()) {
	p.submitLock.Lock()
	defer p.submitLock.Unlock()

	neededWorker := 0

	if p.latestUsedWorker != len(p.workers)-1 {
		neededWorker = p.latestUsedWorker + 1
	}

	w := p.workers[neededWorker]
	p.latestUsedWorker = neededWorker

	w.tasksChannel <- f
}

func (p *Pool) Close() {
	for _, w := range p.workers {
		w.closeChannel <- true
	}
}

func NewPool(size int) *Pool {
	readyWg := &sync.WaitGroup{}

	pool := &Pool{
		workers:    make([]*Worker, size),
		submitLock: sync.Mutex{},
	}

	for i := 0; i < size; i++ {
		readyWg.Add(1)

		w := NewWorker(readyWg)
		go w.Run()
		pool.workers[i] = w
	}

	readyWg.Wait()

	return pool
}
