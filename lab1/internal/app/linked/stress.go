package linked

import (
	"fmt"
	"strings"

	"github.com/AndreSS-ntp/univ_algs/lab1/internal/domain"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/pkg/trycatch"
)

// MemoryOverflowError сообщает о том, что во время стресс-теста закончилась память.
type MemoryOverflowError struct {
	cause error
}

func (e *MemoryOverflowError) Error() string {
	return "во время стресс-теста закончилась память"
}

// Unwrap позволяет извлечь исходную ошибку, чтобы сохранить сообщение рантайма.
func (e *MemoryOverflowError) Unwrap() error {
	return e.cause
}

// Is позволяет использовать errors.Is для проверки на MemoryOverflowError.
func (e *MemoryOverflowError) Is(target error) bool {
	_, ok := target.(*MemoryOverflowError)
	return ok
}

// TryFillUntilOOM добавляет элементы в очередь до тех пор, пока не произойдёт паника
// из-за нехватки памяти. Паника перехватывается и возвращается в виде ошибки,
// чтобы основная программа могла продолжить работу.
func TryFillUntilOOM(target *LinkedQueue) error {
	err := trycatch.Do(func() {
		runAllocationLoop(target)
	})
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "out of memory") {
		return &MemoryOverflowError{cause: err}
	}
	return err
}

// runAllocationLoop бесконечно добавляет элементы в переданную очередь,
// чтобы быстро заполнить всю доступную память.
func runAllocationLoop(target *LinkedQueue) {
	counter := 0
	for {
		counter++
		code := fmt.Sprintf("ST%04d", counter%10000)
		part := domain.NewPart(code, 1)
		target.Enqueue(part)
	}
}
