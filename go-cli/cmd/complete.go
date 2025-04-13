package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete <taskid>",
	Short: "Mark a task as complete",
	Long:  `Mark the task with the specified ID as complete.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		fmt.Fprintf(os.Stdout, "Marked task %s as complete\n", taskID)
		// TODO: Implement task completion logic
	},
}

func init() {
	// Add flags specific to this command
}
