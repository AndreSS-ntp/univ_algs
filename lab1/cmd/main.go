package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/AndreSS-ntp/univ_algs/lab1/internal/app/linear"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/app/linked"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/domain"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/pkg"
)

//
// ===== Основной модуль (диалог с пользователем, логика моделирования) =====
//

type appState struct {
	in            *bufio.Reader
	q             domain.Queue
	impl          int
	currentTime   int
	exitRequested bool
	snapshot      mirrorQueue
}

type mirrorNode struct {
	val  domain.ElType
	next *mirrorNode
}

type mirrorQueue struct {
	head *mirrorNode
	tail *mirrorNode
	size int
}

func (mq *mirrorQueue) reset() {
	mq.head, mq.tail = nil, nil
	mq.size = 0
}

func (mq *mirrorQueue) enqueue(v domain.ElType) {
	n := &mirrorNode{val: v}
	if mq.tail == nil {
		mq.head, mq.tail = n, n
	} else {
		mq.tail.next = n
		mq.tail = n
	}
	mq.size++
}

func (mq *mirrorQueue) dequeue() {
	if mq.head == nil {
		return
	}
	mq.head = mq.head.next
	if mq.head == nil {
		mq.tail = nil
	}
	mq.size--
	if mq.size < 0 {
		mq.size = 0
	}
}

func (mq *mirrorQueue) toLinkedQueue() *linked.LinkedQueue {
	cloned := &linked.LinkedQueue{}
	cloned.Init()
	for n := mq.head; n != nil; n = n.next {
		cloned.Enqueue(n.val)
	}
	return cloned
}

func (mq *mirrorQueue) len() int { return mq.size }

func newAppState(in *bufio.Reader) *appState {
	return &appState{in: in}
}

func main() {
	state := newAppState(bufio.NewReader(os.Stdin))
	for {
		if restart := state.run(); restart {
			continue
		}
		if state.exitRequested {
			break
		}
	}
}

func (s *appState) run() (restart bool) {
	defer func() {
		if r := recover(); r != nil {
			if isOOM(r) && s.impl == 2 {
				fmt.Println("\n[!] Недостаточно памяти. Очередь восстановлена до последнего безопасного состояния.")
				s.restoreLinkedQueue()
				fmt.Printf("[!] В очереди %d элементов. Удалите часть элементов и продолжайте работу.\n", s.snapshot.len())
				restart = true
				return
			}
			panic(r)
		}
	}()

	if s.q == nil {
		s.chooseImplementation()
	}

	s.menuLoop()
	return false
}

func (s *appState) chooseImplementation() {
	fmt.Println("Выберите реализацию очереди:")
	fmt.Println("1 — последовательная память (кольцевая, максимум 5 деталей)")
	fmt.Println("2 — связная память (список, без ограничения)")
	impl := pkg.ReadInt(s.in, "Ваш выбор (1/2): ")
	s.applyImplementation(impl)
}

func (s *appState) applyImplementation(impl int) {
	s.impl = impl
	if impl == 2 {
		s.q = &linked.LinkedQueue{}
	} else {
		s.q = &linear.ArrayQueue{}
		s.impl = 1
	}
	s.q.Init()
	s.currentTime = 0
	s.snapshot.reset()
}

func (s *appState) menuLoop() {
	for {
		fmt.Println("\nМеню:")
		fmt.Println("1 — Поставить деталь на обработку (enqueue)")
		fmt.Println("2 — Перейти к следующему моменту модельного времени")
		fmt.Println("3 — Снять текущую деталь с обработки (отказ установки)")
		fmt.Println("4 — Показать очередь")
		fmt.Println("5 — Сбросить процесс (инициализация)")
		fmt.Println("6 — Сменить реализацию очереди")
		fmt.Println("7 — Заполнить связную очередь до исчерпания памяти")
		fmt.Println("0 — Выход")
		choice := pkg.ReadInt(s.in, "Ваш выбор: ")

		switch choice {
		case 1:
			code := pkg.NormalizeCode(pkg.ReadLine(s.in, "Код детали (4 символа, без пробелов): "))
			t := pkg.ReadPositiveInt(s.in, "Время обработки (целое > 0): ")
			p := domain.NewPart(code, t)
			if !s.q.Enqueue(p) {
				fmt.Println("Очередь переполнена — нельзя добавить деталь.")
			} else {
				fmt.Printf("Деталь %s поставлена в очередь. Требуемое время: %d\n", p.Code, p.StartTime)
				s.recordEnqueue(p)
			}

		case 2:
			s.currentTime++
			if s.q.Empty() {
				fmt.Printf("[t=%d] Очередь пуста, установка простаивает.\n", s.currentTime)
				continue
			}
			if first, ok := s.q.Front(); ok {
				first.Time--
				if first.Time <= 0 {
					done, _ := s.q.Dequeue()
					fmt.Printf("[t=%d] Обработка детали %s завершена (время=%d). Деталь исключена из очереди.\n",
						s.currentTime, done.Code, done.StartTime)
					s.recordDequeue()
				} else {
					fmt.Printf("[t=%d] Идёт обработка детали %s, осталось %d.\n",
						s.currentTime, first.Code, first.Time)
				}
			}

		case 3:
			if s.q.Empty() {
				fmt.Println("Очередь пуста — снимать нечего.")
			} else {
				p, _ := s.q.Dequeue()
				fmt.Printf("Деталь %s снята с обработки досрочно (оставалось %d).\n", p.Code, p.Time)
				s.recordDequeue()
			}

		case 4:
			pkg.ShowQueue(s.q, s.currentTime)

		case 5:
			s.q.Init()
			s.currentTime = 0
			s.snapshot.reset()
			fmt.Println("Моделирование сброшено: очередь и модельное время обнулены.")

		case 6:
			impl := pkg.ReadInt(s.in, "Новая реализация (1 — кольцевая, 2 — связная): ")
			s.applyImplementation(impl)
			fmt.Println("Реализация переключена. Очередь и время сброшены.")

		case 7:
			s.fillLinkedQueue()

		case 0:
			fmt.Println("Завершение работы.")
			s.exitRequested = true
			return

		default:
			fmt.Println("Нет такого пункта меню.")
		}
	}
}

func (s *appState) recordEnqueue(p domain.ElType) {
	if s.impl == 2 {
		s.snapshot.enqueue(p)
	}
}

func (s *appState) recordDequeue() {
	if s.impl != 2 || s.snapshot.len() == 0 {
		return
	}
	s.snapshot.dequeue()
}

func (s *appState) restoreLinkedQueue() {
	s.q = s.snapshot.toLinkedQueue()
}

func (s *appState) fillLinkedQueue() {
	if s.impl != 2 {
		fmt.Println("Опция доступна только для связной реализации очереди.")
		return
	}
	fmt.Println("Автоматическое заполнение очереди. Добавление будет продолжаться до исчерпания памяти...")
	count := s.snapshot.len()
	for {
		if free, err := freeMemory(); err == nil && free < minFreeMemory {
			panic(errSyntheticOOM)
		}
		code := fmt.Sprintf("%04X", count%65536)
		part := domain.NewPart(code, 1)
		s.q.Enqueue(part)
		s.recordEnqueue(part)
		count++
		if count%10000 == 0 {
			fmt.Printf("Добавлено %d деталей.\n", count)
		}
	}
}

const minFreeMemory = 32 << 20 // 32 МБ оставляем системе

var errSyntheticOOM = errors.New("out of memory (detected)")

func freeMemory() (uint64, error) {
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		return 0, err
	}
	unit := uint64(info.Unit)
	if unit == 0 {
		unit = 1
	}
	return uint64(info.Freeram) * unit, nil
}

func isOOM(r interface{}) bool {
	switch v := r.(type) {
	case error:
		return strings.Contains(v.Error(), "out of memory")
	case string:
		return strings.Contains(v, "out of memory")
	default:
		return strings.Contains(fmt.Sprint(v), "out of memory")
	}
}
