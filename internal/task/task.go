package task

type filepath = string

type Task struct {
	Name      string
	Status    TaskStatus
	Describe  string
	Target    filepath
	Dependent []filepath
}

type TaskStatus string

const (
	Pending   TaskStatus = "Pending"
	Running   TaskStatus = "Running"
	Completed TaskStatus = "Completed"
)
