package async

import (
	"runtime"
	"sync"
)

type Task func() error

var DefaultPoolSize = PoolSize(8)

func PoolSize(size int) int {
	return size * runtime.NumCPU()
}

type Pool struct {
	Tasks       []Task
	errors      []error
	mutex       *sync.Mutex
	concurrency int
	tasksChan   chan Task
	wg          sync.WaitGroup
}

// NewPool initializes a new pool with the given tasks and
// at the given concurrency.
func NewPool(concurrency int) *Pool {
	return &Pool{
		Tasks:       make([]Task, 0),
		errors:      make([]error, 0),
		concurrency: concurrency,
		tasksChan:   make(chan Task),
		mutex:       &sync.Mutex{},
	}
}

func (p *Pool) AddTask(task Task) {
	p.mutex.Lock()
	p.Tasks = append(p.Tasks, task)
	p.mutex.Unlock()
}

// Run runs all work within the pool and blocks until it's
// finished.
func (p *Pool) Run() {
	for i := 0; i < p.concurrency; i++ {
		go p.work()
	}

	p.wg.Add(len(p.Tasks))
	for _, task := range p.Tasks {
		p.tasksChan <- task
	}

	// all workers return
	close(p.tasksChan)
	p.wg.Wait()
}

// The work loop for any single goroutine.
func (p *Pool) work() {
	for task := range p.tasksChan {
		err := task()
		if err != nil {
			p.mutex.Lock()
			p.errors = append(p.errors, err)
			p.mutex.Unlock()
		}
		p.wg.Done()
	}
}
