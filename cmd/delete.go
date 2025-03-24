package cmd

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	t "todolist/tasks"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use: "delete id",
	Short: "Delete a task",
	Args: cobra.ExactArgs(1),
	RunE: deleteTask,
}

func deleteTask(cmd *cobra.Command, args []string) error {
	tasks, err := t.GetTasks()
	if err != nil {
		return fmt.Errorf("deleteTasks: %w", err)
	}

	targetId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid id: must be a number")
	}

	index := -1
	for i := range tasks {
		if tasks[i].Id == targetId {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("task with id %d not found", targetId)
	}

	tasks = slices.Delete(tasks, index, index+1)

	return t.WriteTasks(tasks)
}