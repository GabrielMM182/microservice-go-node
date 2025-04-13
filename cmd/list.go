package cmd

import (
	"fmt"
	"os"
	"tasks/internal/storage"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

var all bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long:  `List all uncompleted tasks. Use -a or --all to list all tasks including completed ones.`,
	Run: func(cmd *cobra.Command, args []string) {
		if all {
			fmt.Fprintln(os.Stdout, "Listing all tasks...")
		} else {
			fmt.Fprintln(os.Stdout, "Listing uncompleted tasks...")
		}

		// Open and lock the data file
		f, err := storage.LoadFile(storage.DataFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		defer storage.CloseFile(f)

		// Read tasks from file
		tasks, err := storage.ReadTasks(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		// Filter and display tasks
		for _, t := range tasks {
			if !all && t.IsComplete() {
				continue
			}

			status := "[ ]"
			if t.IsComplete() {
				status = "[✓]"
			}

			// Format creation time as relative time
		createdAgo := timediff.TimeDiff(t.CreatedAt)

			// Format due date if present
		dueStr := ""
			if t.DueDate != nil {
				dueStr = fmt.Sprintf(" (due: %s)", t.DueDate.Format("2006-01-02"))
			}

			// Format completed time if present
		completedStr := ""
			if t.IsComplete() {
				completedStr = fmt.Sprintf(" - completed %s", timediff.TimeDiff(*t.CompletedAt))
			}

			fmt.Fprintf(os.Stdout, "%s %d. %s%s (created %s)%s\n",
				status, t.ID, t.Description, dueStr, createdAgo, completedStr)
		}
	},
}

func init() {
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "List all tasks including completed ones")
}
