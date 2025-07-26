package executor

import (
	"ac/internal/command"
	"ac/internal/task"
)

type TaskManager interface {
	Add(task task.Task) error
	List() ([]task.Task, error)
}

type executor struct {
	taskManager TaskManager
}

func NewExecutor(tm TaskManager) *executor {
	return &executor{
		taskManager: tm,
	}
}

func (e *executor) Execute(command.Cmd) error {
	_, err := e.taskManager.List()
	if err != nil {
		return err
	}
	return nil
}
