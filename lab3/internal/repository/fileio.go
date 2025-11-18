package repository

import (
	"bufio"
	"fmt"
	"github.com/AndreSS-ntp/univ_algs/lab3/internal/pkg"
	"os"
	"strings"
)

// LoadFromFile загружает граф из файла с матрицей смежности.
// Формат файла:
// n
// v1 v2 ... vn
// матрица n x n (0/1)
func LoadFromFile(path string) (*pkg.Graph, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var n int
	if _, err := fmt.Fscan(file, &n); err != nil {
		return nil, fmt.Errorf("не удалось прочитать количество вершин: %v", err)
	}

	g := pkg.NewGraph()
	labels := make([]rune, n)

	// читаем метки вершин
	for i := 0; i < n; i++ {
		var s string
		if _, err := fmt.Fscan(file, &s); err != nil {
			return nil, fmt.Errorf("не удалось прочитать метку вершины: %v", err)
		}
		if s == "" {
			return nil, fmt.Errorf("пустая метка вершины")
		}
		r := []rune(s)[0]
		labels[i] = r
		g.AddVertex(r)
	}

	// читаем матрицу смежности
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var val int
			if _, err := fmt.Fscan(file, &val); err != nil {
				return nil, fmt.Errorf("ошибка чтения матрицы смежности: %v", err)
			}
			if val != 0 {
				if err := g.AddEdge(labels[i], labels[j]); err != nil {
					return nil, err
				}
			}
		}
	}

	return g, nil
}

// readRune читает одну букву (rune) из stdin.
func ReadRune(reader *bufio.Reader) (rune, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, fmt.Errorf("пустой ввод")
	}
	return []rune(line)[0], nil
}
