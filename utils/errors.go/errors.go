package utils

import (
	"errors"
	"fmt"
)

// ErrTaskNotFound — кастомная структура для ошибки "Задача не найдена".
type ErrTaskNotFound struct {
	ID int
}

// Error реализует интерфейс error. Теперь структуру можно возвращать как ошибку.
func (e *ErrTaskNotFound) Error() string {
	return fmt.Sprintf("Задача с ID %d не найдена", e.ID)
}

// IsErrTaskNotFound — помощник для проверки типа ошибки.
func IsErrTaskNotFound(err error) bool {
	var target *ErrTaskNotFound
	return errors.As(err, &target)
}
