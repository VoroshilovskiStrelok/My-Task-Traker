package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateTimestamp обновляет время последнего изменения.
func (t *Task) UpdateTimestamp() {
	t.UpdatedAt = time.Now() // UTC по умолчанию, норм для JSON
}

func LoadTasks() ([]Task, error) {
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, fmt.Errorf("Не удалось прочитать файл tasks.json: %w", err)
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("Не удалось десериализовать JSON: %w", err)
	}

	return tasks, nil
}

// SaveTasks сохраняет задачи в tasks.json с отступами.
func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ") // "", "  " — для pretty-print
	if err != nil {
		return fmt.Errorf("Не удалось организовать выполнение задач: %w", err)
	}
	if err := os.WriteFile("tasks.json", data, 0644); err != nil { // 0644 — права доступа (rw-r--r--)
		return fmt.Errorf("Не удалось записать файл tasks.json: %w", err)
	}
	return nil // Всё ок
}
