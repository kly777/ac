package manager

import (
	"ac/internal/task"
)

type taskManager struct {
	db *task.TaskDB
}

func NewTaskManager() *taskManager {
	return &taskManager{
		db: task.GetDB(),
	}
}

func (m *taskManager) CreateTask(t *task.Task) error {
	return m.db.CreateTask(t)
}

func (m *taskManager) GetTask(taskID string) (*task.Task, error) {
	return m.db.GetTask(taskID)
}

func (m *taskManager) ListTasks() ([]*task.Task, error) {
	return m.db.ListTasks()
}

func (m *taskManager) UpdateTaskStatus(taskID string, status string) error {
	task, err := m.GetTask(taskID)
	if err != nil {
		return err
	}
	task.Status = status
	return m.db.UpdateTask(task)
}
