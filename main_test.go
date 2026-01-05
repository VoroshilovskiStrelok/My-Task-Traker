package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"ty-task-tracker/models"
)

func TestCLIIntegration(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Собираем программу один раз перед тестами
	exePath := filepath.Join(tmpDir, "task-cli.exe")
	buildCmd := exec.Command("go", "build", "-o", exePath, "main.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Не удалось собрать бинарник: %v", err)
	}

	// 2. Тест команды ADD
	t.Run("Команда ADD", func(t *testing.T) {
		// Запускаем уже собранный файл из временной папки
		cmd := exec.Command(exePath, "add", "Интеграционный тест")
		cmd.Dir = tmpDir

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Команда add провалилась: %v, вывод: %s", err, output)
		}

		if !strings.Contains(string(output), "Задача успешно добавлена") {
			t.Errorf("Неверный вывод: %s", output)
		}

		// Проверяем JSON
		jsonPath := filepath.Join(tmpDir, "tasks.json")
		if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
			t.Fatal("Файл tasks.json не создан")
		}
	})

	// 3. Тест команды DELETE
	t.Run("Команда DELETE ошибка 999", func(t *testing.T) {
		cmd := exec.Command(exePath, "delete", "999")
		cmd.Dir = tmpDir

		output, _ := cmd.CombinedOutput()

		if !strings.Contains(string(output), "Задача с ID 999 не найдена") {
			t.Errorf("Ожидалась ошибка про ID 999, но получено: %s", output)
		}
	})
}

func TestAddTaskLogic(t *testing.T) {
	tests := []struct {
		name    string
		initial []models.Task
		desc    string
		wantLen int
		wantID  int
		wantErr bool
	}{
		{
			name:    "Добавление в пустой список",
			initial: []models.Task{},
			desc:    "Первая задача",
			wantLen: 1,
			wantID:  1,
			wantErr: false,
		},
		{
			name: "Добавление к существующим (проверка ID)",
			initial: []models.Task{
				{ID: 10, Description: "Старая"},
			},
			desc:    "Новая",
			wantLen: 2,
			wantID:  11,
			wantErr: false,
		},
		{
			name:    "Пустое описание (ошибка)",
			initial: []models.Task{},
			desc:    "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTasks, gotID, err := addTaskLogic(tt.initial, tt.desc)

			if (err != nil) != tt.wantErr {
				t.Errorf("addTaskLogic() ошибка = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(gotTasks) != tt.wantLen {
					t.Errorf("Длина списка = %d, ожидалось %d", len(gotTasks), tt.wantLen)
				}
				if gotID != tt.wantID {
					t.Errorf("ID новой задачи = %d, ожидалось %d", gotID, tt.wantID)
				}
			}
		})
	}
}
