package task

import "time"

// Task represents a single task in the todo list
type Task struct {
	ID          int       `csv:"ID"`
	Description string    `csv:"Description"`
	CreatedAt   time.Time `csv:"CreatedAt"`
	CompletedAt *time.Time `csv:"CompletedAt"`
	DueDate     *time.Time `csv:"DueDate"`
}

// IsComplete checks if the task is completed
func (t *Task) IsComplete() bool {
	return t.CompletedAt != nil
}
