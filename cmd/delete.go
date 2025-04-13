package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <taskid>",
	Short: "Delete a task",
	Long:  `Delete the task with the specified ID from the task list.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		fmt.Fprintf(os.Stdout, "Deleted task %s\n", taskID)
		// TODO: Implement task deletion logic
	},
}

func init() {
	// Add flags specific to this command
}
