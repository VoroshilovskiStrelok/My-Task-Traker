package main

import (
	"fmt"
	"ty-task-tracker/models"
)

func main() {
	// Оставляем только загрузку
	tasks, err := models.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("You have 0 tasks. Add some!")
		return
	}

	// Оставляем красивый вывод списка
	fmt.Printf("You have %d tasks loaded.\n", len(tasks))
	for i, t := range tasks {
		fmt.Printf("%d. ID: %d | %s [%s] (Created: %s)\n",
			i+1, t.ID, t.Description, t.Status, t.CreatedAt.Format("2006-01-02 15:04"))
	}
}
