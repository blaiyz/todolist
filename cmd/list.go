package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"todolist/tasks"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

var all bool

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&all, "all", "a", false, "include completed tasks in addition to uncompleted tasks")
}

var listCmd = &cobra.Command{
	Use: "list",
	Short: "List all uncompleted tasks",
	RunE: listTasks,
}

func listTasks(cmd *cobra.Command, args []string) error {
	tasks, err := tasks.GetTasks()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 8, 0, 4, ' ', 0)
	defer w.Flush()

	if len(tasks) == 0 {
		fmt.Println("No tasks found!")
	}

	if all {
		fmt.Fprintln(w, "ID\tTask\tCreated\tDone")

		for _, task := range tasks {
			relative := timediff.TimeDiff(task.Created)
			fmt.Fprintf(w, "%d\t%s\t%s\t%t\n", task.Id, task.Description, relative, task.IsComplete)
		}
	} else {
		fmt.Fprintln(w, "ID\tTask\tCreated")

		for _, task := range tasks {
			if task.IsComplete {
				continue
			}

			relative := timediff.TimeDiff(task.Created)
			fmt.Fprintf(w, "%d\t%s\t%s\n", task.Id, task.Description, relative)
		}
	}


	return nil
}