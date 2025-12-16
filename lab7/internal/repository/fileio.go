package graphio

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func stripComments(line string) string {
	if i := strings.Index(line, "#"); i >= 0 {
		line = line[:i]
	}
	if i := strings.Index(line, "//"); i >= 0 {
		line = line[:i]
	}
	return line
}

func splitTokens(line string) []string {
	f := func(r rune) bool {
		return r == ' ' || r == '\t' || r == ',' || r == ';' || r == '|'
	}
	return strings.FieldsFunc(line, f)
}

func isIntToken(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func parseIntRow(tokens []string) ([]int, error) {
	row := make([]int, len(tokens))
	for i, t := range tokens {
		v, err := strconv.Atoi(t)
		if err != nil {
			return nil, fmt.Errorf("не число в матрице: %q", t)
		}
		row[i] = v
	}
	return row, nil
}

func indexToLabel(i int) string {
	i++
	var b []byte
	for i > 0 {
		rem := (i - 1) % 26
		b = append([]byte{byte('A' + rem)}, b...)
		i = (i - 1) / 26
	}
	return string(b)
}

func generateLabels(n int) []string {
	labels := make([]string, n)
	for i := 0; i < n; i++ {
		labels[i] = indexToLabel(i)
	}
	return labels
}

// ParseGraphFile читает граф из файла с матрицей смежности.
// Поддерживает:
// A) первая строка - метки вершин, далее матрица (с меткой строки или без)
// B) только матрица (метки генерируются A,B,C,...)
func ParseGraphFile(path string) ([]string, [][]int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	lines := strings.Split(string(data), "\n")
	var tokenLines [][]string
	for _, raw := range lines {
		line := strings.TrimSpace(stripComments(raw))
		if line == "" {
			continue
		}
		toks := splitTokens(line)
		if len(toks) == 0 {
			continue
		}
		tokenLines = append(tokenLines, toks)
	}
	if len(tokenLines) == 0 {
		return nil, nil, errors.New("файл пустой или не содержит данных")
	}

	first := tokenLines[0]
	firstAllInts := true
	for _, t := range first {
		if !isIntToken(t) {
			firstAllInts = false
			break
		}
	}

	// Вариант B: только матрица
	if firstAllInts {
		n := len(first)
		if len(tokenLines) != n {
			return nil, nil, fmt.Errorf("матрица должна быть квадратной: в первой строке %d элементов, а строк %d", n, len(tokenLines))
		}
		matrix := make([][]int, n)
		for i := 0; i < n; i++ {
			if len(tokenLines[i]) != n {
				return nil, nil, fmt.Errorf("строка %d: ожидалось %d элементов, получено %d", i+1, n, len(tokenLines[i]))
			}
			row, err := parseIntRow(tokenLines[i])
			if err != nil {
				return nil, nil, fmt.Errorf("строка %d: %w", i+1, err)
			}
			matrix[i] = row
		}
		return generateLabels(n), matrix, nil
	}

	// Вариант A: первая строка — метки
	labels := append([]string(nil), first...)
	n := len(labels)

	if len(tokenLines) < 1+n {
		return nil, nil, fmt.Errorf("нужно минимум %d строк матрицы после строки меток, найдено %d", n, len(tokenLines)-1)
	}

	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		toks := tokenLines[i+1]

		// допускаем: "A 0 1 0" или "0 1 0"
		if len(toks) == n+1 && !isIntToken(toks[0]) {
			toks = toks[1:]
		}
		if len(toks) != n {
			return nil, nil, fmt.Errorf("строка %d матрицы: ожидалось %d чисел, получено %d", i+1, n, len(toks))
		}
		row, err := parseIntRow(toks)
		if err != nil {
			return nil, nil, fmt.Errorf("строка %d матрицы: %w", i+1, err)
		}
		matrix[i] = row
	}

	return labels, matrix, nil
}
