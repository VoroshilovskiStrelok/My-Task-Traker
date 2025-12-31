package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func main() {
	t := Task{
		ID:          1,
		Description: "Test task",
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	data, _ := json.MarshalIndent(t, "", "  ")
	fmt.Println(string(data))

	var loaded Task
	json.Unmarshal(data, &loaded)
	fmt.Printf("Loaded: %+v\n", loaded)
}
