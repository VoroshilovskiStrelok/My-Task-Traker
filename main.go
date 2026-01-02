package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"ty-task-tracker/models"
)

func main() {
	// Флаг помощи
	help := flag.Bool("help", false, "Показать справку по использованию")
	flag.Parse()

	if *help {
		printUsage()
		return
	}

	// Получаем аргументы после флагов
	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	cmd := args[0] // Первый аргумент - команда

	switch cmd {
	case "add":
		if len(args) < 2 {
			fmt.Println("Ошибка: Использование: task-cli add <описание>")
			os.Exit(1)
		}
		addTask(args[1])
	case "list":
		var filter string
		if len(args) > 1 {
			filter = strings.ToLower(args[1])
			// Валидация фильтра статуса
			valid := map[string]bool{"todo": true, "in-progress": true, "done": true}
			if !valid[filter] {
				fmt.Printf("Ошибка: Неверный фильтр '%s'. Доступные: todo, in-progress, done\n", filter)
				os.Exit(1)
			}
		}
		listTasks(filter)

	default:
		fmt.Printf("Ошибка: Неизвестная команда '%s'. Используйте --help для справки.\n", cmd)
		os.Exit(1)
	}
}

// addTask — добавляет новую задачу (Загрузить + добавить + Сохранить).
func addTask(desc string) {
	tasks, err := models.LoadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		os.Exit(1)
	}

	newID := len(tasks) + 1
	newTask := models.Task{
		ID:          newID,
		Description: desc,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// Обновит UpdatedAt
	newTask.UpdateTimestamp()

	tasks = append(tasks, newTask)

	if err := models.SaveTasks(tasks); err != nil {
		fmt.Printf("Ошибка сохранения задачи: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Задача успешно добавлена (ID: %d, статус: todo)\n", newID)
}

// listTasks — загружает и печатает все задачи.

func listTasks(filter string) {

	tasks, err := models.LoadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println("Задачи не найдены.")
		return
	}

	var filtered []models.Task
	for _, t := range tasks {
		if filter == "" || t.Status == filter {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) == 0 {
		if filter != "" {
			fmt.Printf("Задачи со статусом %s не найдены.\n", filter)
		} else {
			fmt.Println("Задачи не найдены.")
		}
		return
	}

	filterName := filter
	if filterName == "" {
		filterName = "все"
	}

	fmt.Printf("Ваши  %s задачи (%d всего):\n", map[string]string{"": "all", "todo": "todo", "in-progress": "in-progress", "done": "done"}[filter], len(filtered))
	for _, t := range filtered {
		fmt.Printf("ID: %d | %s | Status: %s | Created: %s | Updated: %s\n",
			t.ID, t.Description, t.Status,
			t.CreatedAt.Format("2006-01-02 15:04"), // YYYY-MM-DD HH:MM
			t.UpdatedAt.Format("2006-01-02 15:04"))
	}
}

// printUsage — простая справка.
func printUsage() {
	fmt.Println("Использование: task-cli <command> [args]")
	fmt.Println("Команды:")
	fmt.Println("  add <description>             Добавить новую задачу (status: todo)")
	fmt.Println("  list [todo|in-progress|done]  Показать задачи (фильтр опционален)")
	fmt.Println("  --help                        Показать справку")
}
