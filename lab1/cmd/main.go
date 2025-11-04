package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/AndreSS-ntp/univ_algs/lab1/internal/app/linear"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/app/linked"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/domain"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/pkg"
)

//
// ===== Основной модуль (диалог с пользователем, логика моделирования) =====
//

func main() {
	in := bufio.NewReader(os.Stdin)

	// Выбор реализации
	var q domain.Queue
	fmt.Println("Выберите реализацию очереди:")
	fmt.Println("1 — последовательная память (кольцевая, максимум 5 деталей)")
	fmt.Println("2 — связная память (список, без ограничения)")
	impl := pkg.ReadInt(in, "Ваш выбор (1/2): ")
	if impl == 2 {
		q = &linked.LinkedQueue{}
	} else {
		q = &linear.ArrayQueue{}
	}
	q.Init()

	currentTime := 0

	for {
		fmt.Println("\nМеню:")
		fmt.Println("1 — Поставить деталь на обработку (enqueue)")
		fmt.Println("2 — Перейти к следующему моменту модельного времени")
		fmt.Println("3 — Снять текущую деталь с обработки (отказ установки)")
		fmt.Println("4 — Показать очередь")
		fmt.Println("5 — Сбросить процесс (инициализация)")
		fmt.Println("6 — Сменить реализацию очереди")
		fmt.Println("7 — Стресс-тест: добавить детали до исчерпания памяти")
		fmt.Println("0 — Выход")
		choice := pkg.ReadInt(in, "Ваш выбор: ")

		switch choice {
		case 1:
			code := pkg.NormalizeCode(pkg.ReadLine(in, "Код детали (4 символа, без пробелов): "))
			t := pkg.ReadPositiveInt(in, "Время обработки (целое > 0): ")
			p := domain.NewPart(code, t)
			if !q.Enqueue(p) {
				fmt.Println("Очередь переполнена — нельзя добавить деталь.")
			} else {
				fmt.Printf("Деталь %s поставлена в очередь. Требуемое время: %d\n", p.Code, p.StartTime)
			}

		case 2:
			currentTime++
			if q.Empty() {
				fmt.Printf("[t=%d] Очередь пуста, установка простаивает.\n", currentTime)
				continue
			}
			if first, ok := q.Front(); ok {
				first.Time--
				if first.Time <= 0 {
					done, _ := q.Dequeue()
					fmt.Printf("[t=%d] Обработка детали %s завершена (время=%d). Деталь исключена из очереди.\n",
						currentTime, done.Code, done.StartTime)
				} else {
					fmt.Printf("[t=%d] Идёт обработка детали %s, осталось %d.\n",
						currentTime, first.Code, first.Time)
				}
			}

		case 3:
			if q.Empty() {
				fmt.Println("Очередь пуста — снимать нечего.")
			} else {
				p, _ := q.Dequeue()
				fmt.Printf("Деталь %s снята с обработки досрочно (оставалось %d).\n", p.Code, p.Time)
			}

		case 4:
			pkg.ShowQueue(q, currentTime)

		case 5:
			q.Init()
			currentTime = 0
			fmt.Println("Моделирование сброшено: очередь и модельное время обнулены.")

		case 6:
			impl = pkg.ReadInt(in, "Новая реализация (1 — кольцевая, 2 — связная): ")
			if impl == 2 {
				q = &linked.LinkedQueue{}
			} else {
				q = &linear.ArrayQueue{}
			}
			q.Init()
			currentTime = 0
			fmt.Println("Реализация переключена. Очередь и время сброшены.")

		case 7:
			linkedQueue, ok := q.(*linked.LinkedQueue)
			if !ok {
				fmt.Println("Стресс-тест доступен только для связной реализации очереди.")
				continue
			}

			fmt.Println("Запуск стресс-теста: добавляем детали, пока не закончится память...")
			if err := linked.TryFillUntilOOM(linkedQueue); err != nil {
				var memErr *linked.MemoryOverflowError
				if errors.As(err, &memErr) {
					fmt.Println("Ошибка: закончилась память.")
					fmt.Printf("В очереди сейчас %d элементов. Удалите хотя бы одну деталь, чтобы освободить память и попробовать снова.\n",
						linkedQueue.Len())
				} else {
					fmt.Printf("Стресс-тест завершился с ошибкой: %v\n", err)
				}
				continue
			}

			fmt.Println("Стресс-тест завершён без ошибок.")

		case 0:
			fmt.Println("Завершение работы.")
			return

		default:
			fmt.Println("Нет такого пункта меню.")
		}
	}
}
