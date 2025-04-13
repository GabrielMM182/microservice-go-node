package cmd

import (
	"fmt"
	"os"
	"strings"
	"tasks/internal/storage"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <description>",
	Short: "Add a new task",
	Long:  `Add a new task with the provided description to the task list.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		description := strings.Join(args, " ")

		// Open and lock the data file
		f, err := storage.LoadFile(storage.DataFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		defer storage.CloseFile(f)

		// Read existing tasks
		tasks, err := storage.ReadTasks(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading tasks: %v\n", err)
			return
		}

		// Add the new task
		tasks = storage.AddTask(tasks, description, nil) // Assuming no due date for now

		// Write the updated tasks back to the file
		if err := storage.WriteTasks(f, tasks); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing tasks: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Added task: %s\n", description)
	},
}

func init() {
	// Add flags specific to this command if needed (e.g., --due-date)
}
