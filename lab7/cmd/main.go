package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/AndreSS-ntp/univ_algs/lab7/internal/pkg"
	"github.com/AndreSS-ntp/univ_algs/lab7/internal/repository"
)

func readLine(r *bufio.Reader) (string, error) {
	s, err := r.ReadString('\n')
	if err != nil && len(s) == 0 {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

func printHelp() {
	fmt.Println("Меню:")
	fmt.Println(" 1) Загрузить граф из файла")
	fmt.Println(" 2) Показать вершины")
	fmt.Println(" 3) Показать матрицу смежности")
	fmt.Println(" 4) Выбрать стартовую вершину")
	fmt.Println(" 5) DFS (обход в глубину) от стартовой вершины")
	fmt.Println(" 6) BFS (обход по уровням) от стартовой вершины")
	fmt.Println(" 0) Выход")
}

func printMatrix(g *graph.Graph) {
	n := len(g.Labels)
	fmt.Println("Матрица смежности:")
	fmt.Print("    ")
	for _, lab := range g.Labels {
		fmt.Printf("%3s", lab)
	}
	fmt.Println()
	for i := 0; i < n; i++ {
		fmt.Printf("%3s ", g.Labels[i])
		for j := 0; j < n; j++ {
			v := 0
			if g.Adj[i][j] != 0 {
				v = 1
			}
			fmt.Printf("%3d", v)
		}
		fmt.Println()
	}
}

func main() {
	var fileFlag string
	flag.StringVar(&fileFlag, "file", "", "путь к файлу с матрицей смежности")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	var g *graph.Graph
	start := -1

	loadFromFile := func(path string) {
		labels, matrix, err := graphio.ParseGraphFile(path)
		if err != nil {
			fmt.Println("Ошибка чтения файла:", err)
			return
		}
		gr, err := graph.New(labels, matrix)
		if err != nil {
			fmt.Println("Ошибка построения графа:", err)
			return
		}
		g = gr
		start = 0
		fmt.Printf("Граф загружен. Вершин: %d. Стартовая по умолчанию: %s\n", len(g.Labels), g.Labels[start])
	}

	if fileFlag != "" {
		loadFromFile(fileFlag)
	}

	for {
		fmt.Println()
		printHelp()
		fmt.Print("Выбор: ")
		choice, err := readLine(reader)
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			return
		}

		switch choice {
		case "1":
			fmt.Print("Введите путь к файлу: ")
			path, _ := readLine(reader)
			if path == "" {
				fmt.Println("Путь пустой.")
				continue
			}
			loadFromFile(path)

		case "2":
			if g == nil {
				fmt.Println("Сначала загрузите граф (пункт 1).")
				continue
			}
			fmt.Println("Вершины:", strings.Join(g.Labels, " "))
			if start >= 0 {
				fmt.Println("Текущая стартовая вершина:", g.Labels[start])
			}

		case "3":
			if g == nil {
				fmt.Println("Сначала загрузите граф (пункт 1).")
				continue
			}
			printMatrix(g)

		case "4":
			if g == nil {
				fmt.Println("Сначала загрузите граф (пункт 1).")
				continue
			}
			fmt.Print("Введите метку стартовой вершины (например A): ")
			lab, _ := readLine(reader)
			if i, ok := g.VertexIndex(lab); ok {
				start = i
				fmt.Println("Стартовая вершина установлена:", g.Labels[start])
			} else {
				fmt.Println("Нет такой вершины:", lab)
			}

		case "5":
			if g == nil || start < 0 {
				fmt.Println("Сначала загрузите граф и выберите стартовую вершину.")
				continue
			}
			ord := g.DFS(start)
			fmt.Println("DFS порядок обхода:", g.FormatOrder(ord))

		case "6":
			if g == nil || start < 0 {
				fmt.Println("Сначала загрузите граф и выберите стартовую вершину.")
				continue
			}
			ord := g.BFS(start)
			fmt.Println("BFS порядок обхода:", g.FormatOrder(ord))

		case "0":
			fmt.Println("Выход.")
			return

		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}
