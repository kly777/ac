package task

import (
	"ac/internal/db"
	"database/sql"
)

type Dependency struct {
	File           string
	RequiredExports []string
	Version        string
}

type Task struct {
	ID           string
	File         string
	Description  string
	Status       string
	Dependencies []Dependency
	Document     string
}

func NewTask(id, file, description string) *Task {
	return &Task{
		ID:          id,
		File:        file,
		Description: description,
		Status:      "pending",
	}
}

func (t *Task) SetDocument(doc string) {
	t.Document = doc
}

// TaskDB 提供任务数据库操作
type TaskDB struct {
	db *sql.DB
}

// GetDB 返回任务数据库实例
func GetDB() *TaskDB {
	return &TaskDB{db: db.DB}
}

// CreateTask 创建新任务
func (tdb *TaskDB) CreateTask(t *Task) error {
	_, err := tdb.db.Exec(`
		INSERT INTO tasks (id, file, description, status, dependencies, document)
		VALUES (?, ?, ?, ?, ?, ?)`,
		t.ID, t.File, t.Description, t.Status, t.Dependencies, t.Document)
	return err
}

// UpdateTask 更新任务状态
func (tdb *TaskDB) UpdateTask(t *Task) error {
	_, err := tdb.db.Exec(`
		UPDATE tasks SET status = ?, document = ?
		WHERE id = ?`,
		t.Status, t.Document, t.ID)
	return err
}

// GetTask 获取任务
func (tdb *TaskDB) GetTask(taskID string) (*Task, error) {
	var t Task
	err := tdb.db.QueryRow("SELECT id, file, description, status, dependencies, document FROM tasks WHERE id = ?", taskID).
		Scan(&t.ID, &t.File, &t.Description, &t.Status, &t.Dependencies, &t.Document)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ListTasks 列出所有任务
func (tdb *TaskDB) ListTasks() ([]*Task, error) {
	rows, err := tdb.db.Query("SELECT id, file, description, status, dependencies, document FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.File, &t.Description, &t.Status, &t.Dependencies, &t.Document); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}