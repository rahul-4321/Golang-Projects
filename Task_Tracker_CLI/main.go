package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `jsson:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

const fileName = "tasks.json"

func loadTasks() ([]Task, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	var tasks []Task

	if err := json.Unmarshal(file, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func addTask(description string) {
	tasks, _ := loadTasks()
	id := len(tasks) + 1
	currentTime := time.Now().Format(time.RFC3339)
	newTask := Task{id, description, "todo", currentTime, currentTime}

	tasks = append(tasks, newTask)

	if err := saveTasks(tasks); err != nil {
		log.Println("Error in saving the task: ", err)
		return
	}
	log.Println("Task added successfully", id)
}

func updateTask(id int, status string) {
	tasks, _ := loadTasks()

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			if err := saveTasks(tasks); err != nil {
				log.Println("Error in updating the task: ", err)
				return
			}
			log.Println("Task updated successfully", id)
		}
	}
	log.Println("Task not found with ID: ", id)
}

func deleteTask(id int) {
	tasks, _ := loadTasks()
	tasksAfterDelete := []Task{}
	found := false
	for i := range tasks {
		if tasks[i].ID != id {
			tasksAfterDelete = append(tasksAfterDelete, tasks[i])
		} else {
			found = true
		}
	}

	if !found {
		log.Println("Task not found with ID: ", id)
		return
	}

	saveTasks(tasksAfterDelete)
	log.Println("Task deleted successfully", id)
}

func markTaskStatus(id int, status string) {
	tasks, _ := loadTasks()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			if err := saveTasks(tasks); err != nil {
				log.Println("Error in marking the task status: ", err)
				return
			}
		}
	}
	log.Printf("Task %d marked %s successfully", id, status)
}

func listTasks(status string) {
	tasks, _ := loadTasks()

	if len(tasks) == 0 {
		log.Println("No tasks found")
		return
	}
	fmt.Println("TaskId	 | 	Description	  |  Status	 |   CreatedAt	|	UpdatedAt")
	for _, task := range tasks {
		if task.Status == status || status == " " {
			fmt.Printf("%d	  |    %s  | %s 	 | %s	| 	%s", task.ID, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
		}
	}
}

func printHelp() {
	helpText := `
Usage: ./task-tracker <command> [arguments]

Commands:
  add "description"           Add a new task
  update <id> "description"   Update a task's description
  delete <id>                 Delete a task
  mark-todo <id>              Mark a task as todo
  mark-in-progress <id>       Mark a task as in progress
  mark-done <id>              Mark a task as done
  list                        List all tasks
  list [status]               List tasks filtered by status (todo, in-progress, done)
  --help                      Show this help message
`
	fmt.Println(helpText)
}

func main() {

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	if os.Args[1] == "--help" {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {

	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a description for the task.")
			return
		}
		description := strings.Join(os.Args[2:], " ")
		addTask(description)

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Please provide an ID and a new description for the task.")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		description := strings.Join(os.Args[3:], " ")
		updateTask(id, description)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Please provide an ID of the task to delete.")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		deleteTask(id)

	case "mark-todo":
		if len(os.Args) < 3 {
			fmt.Println("Please provide an ID of the task to mark as todo.")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		markTaskStatus(id, "todo")

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Please provide an ID of the task to mark as in progress.")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		markTaskStatus(id, "in-progress")

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Please provide an ID of the task to mark as done.")
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		markTaskStatus(id, "done")

	case "list":
		if len(os.Args) == 2 {
			listTasks("")
		} else if len(os.Args) == 3 {
			status := os.Args[2]
			listTasks(status)
		}
	}
}
