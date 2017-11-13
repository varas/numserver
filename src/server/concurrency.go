package server

import "sync"

type concurrencyManager struct {
	sync.Mutex
	tasks    int
	maxTasks int
	wg       sync.WaitGroup
}

func newConcurrencyManager(maxConcurrentTasks int) *concurrencyManager {
	return &concurrencyManager{
		tasks:    0,
		maxTasks: maxConcurrentTasks,
		wg:       sync.WaitGroup{},
	}
}

func (c *concurrencyManager) AddTaskOrWait() {
	if c.tasks >= c.maxTasks {
		c.wg.Wait()
	}

	c.Lock()
	c.wg.Add(1)
	c.tasks++
	c.Unlock()
}

func (c *concurrencyManager) FinishTask() {
	c.Lock()
	c.wg.Done()
	c.tasks--
	c.Unlock()
}
