package storage

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"
	"tasks/internal/task"
	"time"
)

// Storage manages the task data file
const (
	DataFilePath = "tasks.csv"
)

// LoadFile opens and locks the data file for reading and writing
func LoadFile(filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %v", err)
	}

	// Exclusive lock obtained on the file descriptor
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("failed to lock file: %v", err)
	}

	return f, nil
}

// CloseFile unlocks and closes the file
func CloseFile(f *os.File) error {
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		return fmt.Errorf("failed to unlock file: %v", err)
	}
	return f.Close()
}

// ReadTasks reads tasks from the CSV file
func ReadTasks(f *os.File) ([]task.Task, error) {
	f.Seek(0, 0)
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header
	_, err := reader.Read()
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error reading CSV header: %v", err)
	}

	var tasks []task.Task
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %v", err)
		}

		if len(record) < 3 {
			continue // Skip invalid records
		}

		id, _ := strconv.Atoi(record[0])
		createdAt, _ := time.Parse(time.RFC3339, record[2])

		var completedAt *time.Time
		if record[3] != "" {
			t, _ := time.Parse(time.RFC3339, record[3])
			completedAt = &t
		}

		var dueDate *time.Time
		if len(record) > 4 && record[4] != "" {
			t, _ := time.Parse(time.RFC3339, record[4])
			dueDate = &t
		}

		task := task.Task{
			ID:          id,
			Description: record[1],
			CreatedAt:   createdAt,
			CompletedAt: completedAt,
			DueDate:     dueDate,
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// WriteTasks writes tasks to the CSV file
func WriteTasks(f *os.File, tasks []task.Task) error {
	f.Seek(0, 0)
	f.Truncate(0)

	writer := csv.NewWriter(f)

	// Write header
	header := []string{"ID", "Description", "CreatedAt", "CompletedAt", "DueDate"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	// Write tasks
	for _, t := range tasks {
		completedAt := ""
		if t.CompletedAt != nil {
			completedAt = t.CompletedAt.Format(time.RFC3339)
		}

		dueDate := ""
		if t.DueDate != nil {
			dueDate = t.DueDate.Format(time.RFC3339)
		}

		record := []string{
			strconv.Itoa(t.ID),
			t.Description,
			t.CreatedAt.Format(time.RFC3339),
			completedAt,
			dueDate,
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %v", err)
		}
	}

	writer.Flush()
	return writer.Error()
}

// AddTask adds a new task to the list and returns the updated list
func AddTask(tasks []task.Task, description string, dueDate *time.Time) []task.Task {
	newID := 1
	if len(tasks) > 0 {
		newID = tasks[len(tasks)-1].ID + 1
	}
	newTask := task.Task{
		ID:          newID,
		Description: description,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		DueDate:     dueDate,
	}
	tasks = append(tasks, newTask)
	return tasks
}

// CompleteTask marks a task as complete
func CompleteTask(tasks []task.Task, id int) []task.Task {
	for i, t := range tasks {
		if t.ID == id {
			now := time.Now()
			tasks[i].CompletedAt = &now
			break
		}
	}
	return tasks
}

// DeleteTask removes a task from the list
func DeleteTask(tasks []task.Task, id int) []task.Task {
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
	return tasks
}
