package main

import (
	"flag"
	"fmt"
	"os"
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
		desc := args[1]
		addTask(desc)

	case "list":
		listTasks()

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

	fmt.Printf("Задача успешно добавлена (ID: %d)\n", newID)
}

// listTasks — загружает и печатает все задачи.
func listTasks() {
	tasks, err := models.LoadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println("Задачи не найдены.")
		return
	}

	fmt.Println("Ваши задачи:")
	for _, t := range tasks {
		fmt.Printf("ID: %d | %s | Статус: %s | Создано: %s\n",
			t.ID, t.Description, t.Status, t.CreatedAt.Format("2006-01-02 15:04"))
	}
}

// printUsage — простая справка.
func printUsage() {
	fmt.Println("Использование: task-cli <команда> [аргументы]")
	fmt.Println("Команды:")
	fmt.Println("  add <описание>  Добавить новую задачу")
	fmt.Println("  list            Список всех задач")
	fmt.Println("Запустите 'task-cli --help' для получения дополнительной информации.")
}
