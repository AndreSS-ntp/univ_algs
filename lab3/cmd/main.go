package main

import (
	"bufio"
	"fmt"
	"github.com/AndreSS-ntp/univ_algs/lab3/internal/repository"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите имя файла с матрицей смежности: ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)

	g, err := repository.LoadFromFile(path)
	if err != nil {
		fmt.Println("Ошибка загрузки графа:", err)
		return
	}

	fmt.Println("Граф успешно загружен.")

	for {
		fmt.Println("\nМеню:")
		fmt.Println("1 - Показать граф (список смежности)")
		fmt.Println("2 - Добавить вершину")
		fmt.Println("3 - Удалить вершину")
		fmt.Println("4 - Поиск ребра между двумя вершинами")
		fmt.Println("5 - Выход")
		fmt.Print("Ваш выбор: ")

		choiceLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			return
		}
		choice := strings.TrimSpace(choiceLine)

		switch choice {
		case "1":
			g.Print()

		case "2":
			fmt.Print("Введите букву новой вершины: ")
			v, err := repository.ReadRune(reader)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			if !g.AddVertex(v) {
				fmt.Printf("Вершина %c уже существует.\n", v)
				continue
			}
			fmt.Println("Вершина добавлена.")
			fmt.Println("Введите через пробел вершины, в которые идут дуги из новой вершины (пустая строка — без дуг):")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line != "" {
				for _, token := range strings.Fields(line) {
					r := []rune(token)[0]
					if err := g.AddEdge(v, r); err != nil {
						fmt.Println("Ошибка добавления дуги:", err)
					}
				}
			}

		case "3":
			fmt.Print("Введите букву удаляемой вершины: ")
			v, err := repository.ReadRune(reader)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			if g.RemoveVertex(v) {
				fmt.Println("Вершина удалена.")
			} else {
				fmt.Println("Такой вершины нет в графе.")
			}

		case "4":
			fmt.Print("Введите начальную вершину: ")
			from, err := repository.ReadRune(reader)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			fmt.Print("Введите конечную вершину: ")
			to, err := repository.ReadRune(reader)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			if g.HasEdge(from, to) {
				fmt.Printf("Ребро %c -> %c существует: true\n", from, to)
			} else {
				fmt.Printf("Ребро %c -> %c не найдено: false\n", from, to)
			}

		case "5":
			fmt.Println("Выход из программы.")
			return

		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}
