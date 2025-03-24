package tasks

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
)

type Task struct {
	Id          int
	Description string
	Created     time.Time
	IsComplete bool
}

const filename = "tasks.csv"

func initCsv(w *csv.Writer) error {
	headers := []string{"Id", "Description", "Created", "IsComplete"}
	if err := w.Write(headers); err != nil {
		return err
	}
	w.Flush()

	return w.Error()
}

func openFile() (*os.File, error) {
	_, err := os.Stat("tasks.csv")
	notExists := false
	if os.IsNotExist(err) {
		notExists = true
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, err
	}

	if notExists {
		writer := csv.NewWriter(file)
		if err := initCsv(writer); err != nil {
			return nil, fmt.Errorf("openFile: Couldn't initialize csv: %w", err)
		}
		fmt.Println("tasks.csv created successfully.")
	}

	return file, nil
}

func closeFile(f *os.File) error {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return f.Close()
}

func GetTasks() ([]Task, error) {
	file, err := openFile()
	if err != nil {
		return nil, fmt.Errorf("GetTasks: error opening file: %w", err)
	}
	defer closeFile(file)

	var tasks []Task

	reader := csv.NewReader(file)
	// Read header
	reader.Read()
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("GetTasks: error reading csv: %w", err)
	}

	for _, record := range records {
		if len(record) != 4 {
			fmt.Fprintf(os.Stderr, "Invalid record (wrong number of fields): %v\n", record)
			continue
		}

		// Conversions
		id, err := strconv.ParseInt(record[0], 10, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid id in: %v\n", record)
			continue
		}

		unixTime, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't parse unix time: %v", record)
			continue
		}

		complete, err := strconv.ParseBool(record[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't parse boolean: %v", record)
		}

		task := Task{
			Id: int(id),
			Description: record[1],
			Created: time.Unix(unixTime, 0),
			IsComplete: complete,
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func WriteTasks(tasks []Task) error {
	file, err := openFile()
	if err != nil {
		return fmt.Errorf("WriteTasks: Couldn't open file: %w", err)
	}
	defer closeFile(file)

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("WriteTasks: failed to truncate file: %w", err)
	}
	writer := csv.NewWriter(file)
	if err := initCsv(writer); err != nil {
		return fmt.Errorf("WriteTasks: Couldn't initialize csv: %w", err)
	}


	for _, task := range tasks {
		record := []string{
			strconv.FormatInt(int64(task.Id), 10),
			task.Description,
			strconv.FormatInt(int64(task.Created.Unix()), 10),
			strconv.FormatBool(task.IsComplete),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("WriteTasks: Couldn't write record %v: %w", record, err)
		}
	}

	writer.Flush()
	return writer.Error()
}