package taskManager

import (
	"ac/internal/task"
	"sync"
)

type taskManager struct {
	mu    sync.RWMutex
	tasks []task.Task
}

func NewTaskManager() *taskManager {
	return &taskManager{}
}

func (tm *taskManager) Add(task task.Task) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tasks = append(tm.tasks, task)
	return nil
}

func (tm *taskManager) List() ([]task.Task, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.tasks, nil
}
