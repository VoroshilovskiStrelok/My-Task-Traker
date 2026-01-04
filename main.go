package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
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
			// Список разрешенных фильтров
			validFilters := map[string]bool{"todo": true, "in-progress": true, "done": true}

			if !validFilters[filter] {
				fmt.Printf("Ошибка: Неверный фильтр '%s'.\nДоступные: todo, in-progress, done (или оставьте пустым для всех задач).\n", filter)
				os.Exit(1)
			}
		}
		listTasks(filter)

	case "update":
		// Проверяем, что переданы и ID, и новое описание
		if len(args) < 3 {
			fmt.Println("Ошибка: Использование: task-cli update <ID> <новое описание>")
			os.Exit(1)
		}

		idStr := args[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Printf("Ошибка: Неверный ID '%s'. Он должен быть числом.\n", idStr)
			os.Exit(1)
		}

		newDesc := args[2]
		updateTask(id, newDesc)

	case "delete":
		if len(args) < 2 {
			fmt.Println("Ошибка: Использование: task-cli delete <ID>")
			os.Exit(1)
		}
		idStr := args[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Printf("Ошибка: Неверный ID '%s'. Он должен быть числом.\n", idStr)
			os.Exit(1)
		}

		deleteTask(id)

	case "mark-in-progress":
		if len(args) < 2 {
			fmt.Println("Ошибка: Использование: task-cli mark-in-progress <ID>")
			os.Exit(1)
		}
		idStr := args[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Printf("Ошибка: Неверный ID '%s': %v\n", idStr, err)
			os.Exit(1)
		}
		markTask(id, "in-progress")

	case "mark-done":
		if len(args) < 2 {
			fmt.Println("Ошибка: Использование: task-cli mark-done <ID>")
			os.Exit(1)
		}
		idStr := args[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Printf("Ошибка: Неверный ID '%s': %v\n", idStr, err)
			os.Exit(1)
		}
		markTask(id, "done")

	default:
		fmt.Printf("Ошибка: Неизвестная команда '%s'. Используйте --help для справки.\n", cmd)
		os.Exit(1)
	}
}

// addTask — добавляет новую задачу (Загрузить + добавить + Сохранить).
func addTask(desc string) {
	desc = strings.TrimSpace(desc) // Чистим пробелы
	if desc == "" {
		fmt.Println("Ошибка: Описание задачи не может быть пустым.")
		os.Exit(1)
	}
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

// updateTask — находит задачу по ID и меняет её описание.
func updateTask(id int, newDesc string) {
	// Валидация: убираем лишние пробелы и проверяем на пустоту
	newDesc = strings.TrimSpace(newDesc)
	if newDesc == "" {
		fmt.Println("Ошибка: Описание задачи не может быть пустым.")
		os.Exit(1)
	}

	tasks, err := models.LoadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		os.Exit(1)
	}

	// Используем наш хелпер из Шага 1
	idx := findTaskIndex(tasks, id)
	if idx == -1 {
		fmt.Printf("Ошибка: Задача с ID %d не найдена.\n", id)
		os.Exit(1)
	}

	// Обновляем данные
	tasks[idx].Description = newDesc
	// Вызываем метод модели для обновления времени
	tasks[idx].UpdateTimestamp()

	// Сохраняем весь слайс обратно в JSON
	if err := models.SaveTasks(tasks); err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Задача %d успешно обновлена.\n", id)
}

func deleteTask(id int) {
	tasks, err := models.LoadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		os.Exit(1)
	}

	idx := findTaskIndex(tasks, id)
	if idx == -1 {
		fmt.Printf("Ошибка: Задача с ID %d не найдена.\n", id)
		os.Exit(1)
	}

	// Удаление: создаем новый слайс, соединяя части ДО и ПОСЛЕ индекса
	tasks = append(tasks[:idx], tasks[idx+1:]...)

	if err := models.SaveTasks(tasks); err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Задача %d успешно удалена.\n", id)
}

func markTask(id int, status string) {
	tasks, err := models.LoadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		os.Exit(1)
	}

	idx := findTaskIndex(tasks, id)
	if idx == -1 {
		fmt.Printf("Ошибка: Задача с ID %d не найдена.\n", id)
		os.Exit(1)
	}

	// Graceful check: если статус уже такой же, ничего не делаем
	if tasks[idx].Status == status {
		fmt.Printf("Задача %d уже имеет статус '%s'.\n", id, status)
		return
	}

	tasks[idx].Status = status
	tasks[idx].UpdateTimestamp() // Обновляем время изменения

	if err := models.SaveTasks(tasks); err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Задача %d успешно переведена в статус: %s\n", id, status)
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
	// фильтрация
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

	// Сортировка по CreatedAt по возрастанию
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].ID < filtered[j].ID
	})

	// ПЕЧАТЬ ТАБЛИЦЫ
	fmt.Printf("\n%-3s | %-30s | %-12s | %-12s | %-12s\n", "ID", "Description", "Status", "Created", "Updated")
	fmt.Println(strings.Repeat("-", 80)) // Разделительная линия

	for _, t := range filtered {
		statusUpper := strings.ToUpper(t.Status)
		// Формат даты: Месяц/День Часы:Минуты
		created := t.CreatedAt.Format("01/02 15:04")
		updated := t.UpdatedAt.Format("01/02 15:04")

		// %-3d — число, 3 символа, выравнивание влево
		// %-30.30s — строка, 30 символов, обрезается если длиннее
		fmt.Printf("%-3d | %-30.30s | %-12s | %-12s | %-12s\n",
			t.ID, t.Description, statusUpper, created, updated)
	}
	fmt.Printf("\nВсего задач в списке: %d\n", len(filtered))
}

// findTaskIndex ищет индекс задачи по ID. -1 если не найден.
func findTaskIndex(tasks []models.Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

// printUsage — простая справка.
func printUsage() {
	fmt.Println("Task Tracker CLI — Управляй своими задачами из терминала")
	fmt.Println("\nИспользование:")
	fmt.Println("  task-cli <command> [arguments]")
	fmt.Println("\nКоманды:")
	fmt.Println("  add <описание>             Добавить новую задачу")
	fmt.Println("  list                       Показать все задачи")
	fmt.Println("  list <статус>              Фильтр по статусу (todo, in-progress, done)")
	fmt.Println("  update <ID> <описание>     Изменить описание задачи")
	fmt.Println("  delete <ID>                Удалить задачу по ID")
	fmt.Println("  mark-in-progress <ID>      Установить статус 'в процессе'")
	fmt.Println("  mark-done <ID>             Установить статус 'выполнено'")
	fmt.Println("\nПримеры:")
	fmt.Println("  ./task-cli add \"Купить хлеб\"")
	fmt.Println("  ./task-cli list todo")
	fmt.Println("  ./task-cli mark-done 1")
}
