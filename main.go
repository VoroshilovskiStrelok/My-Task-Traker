package main

import (
	"fmt"
	"ty-task-tracker/models"
)

func main() {
	tasks, err := models.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}
	fmt.Printf("You have %d tasks loaded.\n", len(tasks))
	for i, t := range tasks {
		fmt.Printf("Task %d: %s [%s]\n", i+1, t.Description, t.Status)
	}
}
