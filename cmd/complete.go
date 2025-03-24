package cmd

import (
	"errors"
	"fmt"
	"strconv"
	t "todolist/tasks"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completeCmd)
}

var completeCmd = &cobra.Command{
	Use: "complete id",
	Short: "Complete an uncompleted task in the list",
	Args: cobra.ExactArgs(1),
	RunE: completeTask,
}

func completeTask(cmd *cobra.Command, args []string) error {
	tasks, err := t.GetTasks()
	if err != nil {
		return fmt.Errorf("completeTask: %w", err)
	}

	targetId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid id: must be a number")
	}

	index := -1
	for i := range tasks {
		if tasks[i].Id == targetId {
			index = i
			tasks[i].IsComplete = true
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("task with id %d not found", targetId)
	}

	return t.WriteTasks(tasks)
}