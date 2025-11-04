package linked

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/AndreSS-ntp/univ_algs/lab1/internal/domain"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/pkg/trycatch"
)

// StressChildEnv — переменная окружения, сигнализирующая дочернему процессу,
// что он запущен в режиме стресс-теста.
const StressChildEnv = "LAB1_LINKED_QUEUE_STRESS_CHILD"

// MemoryOverflowError сообщает о том, что дочерний процесс был принудительно
// остановлен из-за нехватки памяти.
type MemoryOverflowError struct {
	Log string
}

func (e *MemoryOverflowError) Error() string {
	return "во время стресс-теста закончилась память"
}

// Is позволяет использовать errors.Is для проверки на MemoryOverflowError.
func (e *MemoryOverflowError) Is(target error) bool {
	_, ok := target.(*MemoryOverflowError)
	return ok
}

// TryFillUntilOOM запускает бесконечное добавление элементов в очереди в
// дочернем процессе. Если тот завершается из-за нехватки памяти, ошибка
// перехватывается в try/catch и возвращается вызывающему коду без падения
// основной программы.
func TryFillUntilOOM() error {
	return trycatch.Do(func() {
		if err := runStressChild(); err != nil {
			panic(err)
		}
	})
}

// runStressChild запускает новый экземпляр программы, который будет добавлять
// элементы в очередь до тех пор, пока система не откажет в памяти.
func runStressChild() error {
	cmd := exec.Command(os.Args[0])
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=1", StressChildEnv))
	cmd.Stdout = io.Discard

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log := strings.TrimSpace(stderr.String())
		if log == "" {
			log = err.Error()
		}
		if strings.Contains(log, "out of memory") || strings.Contains(log, "fatal error") {
			return &MemoryOverflowError{Log: log}
		}
		return fmt.Errorf("дочерний процесс стресс-теста завершился с ошибкой: %w: %s", err, log)
	}

	return nil
}

// HandleStressChildMode запускает режим стресс-теста для дочернего процесса.
func HandleStressChildMode() bool {
	if os.Getenv(StressChildEnv) != "1" {
		return false
	}

	runAllocationLoop()
	return true
}

// runAllocationLoop бесконечно добавляет элементы в связную очередь,
// имитируя стресс-тест на исчерпание памяти.
func runAllocationLoop() {
	var q LinkedQueue
	q.Init()

	counter := 0
	for {
		counter++
		code := fmt.Sprintf("ST%04d", counter%10000)
		part := domain.NewPart(code, 1)
		q.Enqueue(part)
		// Лёгкое «шумовое» действие, чтобы компилятор не оптимизировал цикл.
		if counter%500000 == 0 {
			fmt.Fprintln(os.Stderr, "stress child: продолжаем заполнять очередь...")
		}
	}
}

// WrapMemoryError помогает красиво вывести сообщение пользователю.
func WrapMemoryError(err error) error {
	var memErr *MemoryOverflowError
	if errors.As(err, &memErr) {
		if memErr.Log == "" {
			return errors.New("память закончилась во время стресс-теста")
		}
		return fmt.Errorf("память закончилась во время стресс-теста:\n%s", memErr.Log)
	}
	return err
}
