package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSaveTasksFrom(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "test_save.json")

	// Фиксированное время для тестов
	testTime := time.Date(2026, 1, 5, 15, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		tasks []Task
		want  string // Ожидаемый фрагмент в JSON
	}{
		{
			name:  "Пустой список",
			tasks: []Task{},
			want:  "[]",
		},
		{
			name: "Одна задача",
			tasks: []Task{{
				ID:          1,
				Description: "Тест сохранения",
				Status:      "todo",
				CreatedAt:   testTime,
				UpdatedAt:   testTime,
			}},
			want: `"description": "Тест сохранения"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Вызываем сохранение
			err := SaveTasksFrom(jsonPath, tt.tasks)
			if err != nil {
				t.Fatalf("SaveTasksFrom() ошибка = %v", err)
			}

			// 2. Читаем файл вручную для проверки содержимого
			data, err := os.ReadFile(jsonPath)
			if err != nil {
				t.Fatalf("не удалось прочитать созданный файл: %v", err)
			}

			// 3. Проверяем, содержит ли файл нужную строку
			if !strings.Contains(string(data), tt.want) {
				t.Errorf("Сохраненный JSON не содержит %q. Получено:\n%s", tt.want, string(data))
			}
		})
	}
}

func TestTask_UpdateTimestamp(t *testing.T) {
	t.Run("проверка обновления времени UpdatedAt", func(t *testing.T) {
		// 1. Создаем задачу со старой датой (1 января 2026)
		oldTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		task := Task{
			ID:        1,
			UpdatedAt: oldTime,
		}

		// 2. Вызываем метод обновления
		task.UpdateTimestamp()

		// 3. Проверяем, что время изменилось
		if task.UpdatedAt.Equal(oldTime) {
			t.Error("Ошибка: UpdatedAt не изменился, остался старым")
		}

		// 4. Проверяем, что новое время актуально (не позже, чем 1 минута назад)
		// Это защищает от установки случайных дат в будущем или слишком далеком прошлом
		if task.UpdatedAt.Before(time.Now().UTC().Add(-1 * time.Minute)) {
			t.Error("Ошибка: UpdatedAt установил слишком старое время")
		}
	})
}
