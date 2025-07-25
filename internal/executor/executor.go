package executor

import (
	"ac/internal/db"
	"database/sql"
	"encoding/json"
	"fmt"
)

type TaskExecutor struct {
	db *sql.DB
}

func NewTaskExecutor() *TaskExecutor {
	return &TaskExecutor{
		db: db.DB,
	}
}

func (e *TaskExecutor) ExecuteTasks() error {
	// Fetch all tasks with their dependencies
	rows, err := e.db.Query("SELECT id, dependencies FROM tasks")
	if err != nil {
		return fmt.Errorf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	// Build dependency graph
	type Graph map[string][]string
	type InDegree map[string]int

	graph := make(Graph)
	inDegree := make(InDegree)
	taskIDs := []string{}

	// Initialize graph nodes
	for rows.Next() {
		var id string
		var depsJSON string
		if err := rows.Scan(&id, &depsJSON); err != nil {
			return fmt.Errorf("failed to scan task: %v", err)
		}

		taskIDs = append(taskIDs, id)
		graph[id] = []string{}
		inDegree[id] = 0

		// Parse dependencies
		var deps []string
		if err := json.Unmarshal([]byte(depsJSON), &deps); err == nil {
			for _, dep := range deps {
				if _, exists := graph[dep]; exists {
					graph[dep] = append(graph[dep], id)
					inDegree[id]++
				}
			}
		}
	}

	// Build execution queue using Kahn's algorithm
	queue := []string{}
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	// Process the graph
	executionOrder := []string{}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		executionOrder = append(executionOrder, node)
		visited[node] = true

		for _, neighbor := range graph[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check for cycles
	if len(executionOrder) != len(taskIDs) {
		return fmt.Errorf("cycle detected in task dependencies")
	}

	// Execute tasks in order
	for _, taskID := range executionOrder {
		// Retrieve task details from DB
		var file string
		var description string
		err := e.db.QueryRow("SELECT file, description FROM tasks WHERE id = ?", taskID).
			Scan(&file, &description)
		if err != nil {
			return fmt.Errorf("failed to get task details: %v", err)
		}

		// Get dependencies' documents
		dependencies := make([]string, 0)
		var depsJSON string
		err = e.db.QueryRow("SELECT dependencies FROM tasks WHERE id = ?", taskID).
			Scan(&depsJSON)
		if err == nil {
			json.Unmarshal([]byte(depsJSON), &dependencies)
		}

		// Retrieve document content for dependencies
		var documentContent string
		for _, dep := range dependencies {
			// Get the file path from dependency task
			var depFile string
			e.db.QueryRow("SELECT file FROM tasks WHERE id = ?", dep).
				Scan(&depFile)

			// Get document content
			e.db.QueryRow("SELECT content FROM file_documents WHERE file_path = ?", depFile).
				Scan(&documentContent)

			// Use document content for AI processing
			fmt.Printf("Using document %s for task %s\n", depFile, taskID)
		}

		// Placeholder for AI code generation
		// In a real implementation, this would call the AI API
		// For now, generate mock code and doc
		codeContent := fmt.Sprintf("// Mock code for %s\n// Requirements: %s\nfunc MockFunction() {}\n", file, description)
		docContent := fmt.Sprintf("## %s Documentation\n\nImplementation details for %s.\n\nDependencies:\n", file, file)
		for _, dep := range dependencies {
			var depFile string
			e.db.QueryRow("SELECT file FROM tasks WHERE id = ?", dep).
				Scan(&depFile)
			docContent += fmt.Sprintf("- %s\n", depFile)
		}

		// Store generated code (mock - in real implementation this would be actual storage)
		fmt.Printf("Generated code for %s:\n%s\n", file, codeContent)

		// Store document in SQLite
		_, err = e.db.Exec(`
			INSERT OR REPLACE INTO file_documents (file_path, content)
			VALUES (?, ?)`,
			file, docContent)
		if err != nil {
			return fmt.Errorf("failed to store document: %v", err)
		}

		// Update task status
		_, err = e.db.Exec("UPDATE tasks SET status = 'completed' WHERE id = ?", taskID)
		if err != nil {
			return fmt.Errorf("failed to update task status: %v", err)
		}
	}

	return nil
}
