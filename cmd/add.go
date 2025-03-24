package cmd

import (
	"fmt"
	"time"
	"todolist/tasks"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use: "add description",
	Short: "Add a new task to the list",
	Args: cobra.ExactArgs(1),
	RunE: addTask,
}

func addTask(cmd *cobra.Command, args []string) error {
	t, err := tasks.GetTasks()
	if err != nil {
		return fmt.Errorf("addTask: %w", err)
	}

	var id int
	if id = 0; len(t) > 0 {
		id = t[len(t)-1].Id + 1
	}

	newTask := tasks.Task{
		Id: id,
		Description: args[0],
		Created: time.Now(),
		IsComplete: false,
	}

	t = append(t, newTask)
	return tasks.WriteTasks(t)
}