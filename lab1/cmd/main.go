package main

import (
	"bufio"
	"fmt"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/app/linear"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/app/linked"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/domain"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/pkg"
	"os"
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
		fmt.Println("7 — Заполнить связную очередь до упора по памяти")
		fmt.Println("0 — Выход")
		choice := pkg.ReadInt(in, "Ваш выбор: ")

		switch choice {
		case 1:
			code := pkg.NormalizeCode(pkg.ReadLine(in, "Код детали (4 символа, без пробелов): "))
			t := pkg.ReadPositiveInt(in, "Время обработки (целое > 0): ")
			p := domain.NewPart(code, t)
			if !q.Enqueue(p) {
				fmt.Println("Очередь переполнена — нельзя добавить деталь.")
				if _, ok := q.(*linked.LinkedQueue); ok {
					if free, err := pkg.MemoryHeadroom(); err == nil {
						fmt.Printf("Свободно примерно %.2f МБ памяти.\n", float64(free)/1024.0/1024.0)
					}
				}
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
			lq, ok := q.(*linked.LinkedQueue)
			if !ok {
				fmt.Println("Опция доступна только для связной реализации очереди.")
				break
			}
			fmt.Println("Автозаполнение очереди. Нажмите Ctrl+C для экстренной остановки.")
			added := 0
			for {
				part := domain.NewPart(fmt.Sprintf("A%03d", added%1000), 1)
				if !lq.Enqueue(part) {
					break
				}
				added++
				if added%10000 == 0 {
					fmt.Printf("Добавлено %d деталей...\n", added)
				}
			}
			if free, err := pkg.MemoryHeadroom(); err == nil {
				fmt.Printf("Остановлено после %d деталей. Свободно около %.2f МБ памяти.\n",
					added, float64(free)/1024.0/1024.0)
			} else {
				fmt.Printf("Остановлено после %d деталей из-за ошибки доступа к памяти: %v\n", added, err)
			}
			if limit, limited, err := pkg.MemoryLimit(); err == nil && limited {
				if used, err := pkg.MemoryUsage(); err == nil {
					fmt.Printf("Использовано %.2f МБ из %.2f МБ допустимых.\n",
						float64(used)/1024.0/1024.0, float64(limit)/1024.0/1024.0)
				}
			}

		case 0:
			fmt.Println("Завершение работы.")
			return

		default:
			fmt.Println("Нет такого пункта меню.")
		}
	}
}
