package graph

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Graph struct {
	Labels []string
	Adj    [][]int
	Neigh  [][]int
	Index  map[string]int
}

func New(labels []string, matrix [][]int) (*Graph, error) {
	n := len(labels)
	if n == 0 {
		return nil, errors.New("нет вершин")
	}
	if len(matrix) != n {
		return nil, errors.New("размер матрицы не совпадает с количеством меток")
	}
	for i := range matrix {
		if len(matrix[i]) != n {
			return nil, errors.New("матрица должна быть квадратной")
		}
	}

	idx := make(map[string]int, n)
	for i, lab := range labels {
		key := strings.ToUpper(strings.TrimSpace(lab))
		if key == "" {
			return nil, fmt.Errorf("пустая метка вершины на позиции %d", i)
		}
		if _, ok := idx[key]; ok {
			return nil, fmt.Errorf("дублирующаяся метка вершины: %q", lab)
		}
		idx[key] = i
	}

	g := &Graph{
		Labels: labels,
		Adj:    matrix,
		Index:  idx,
		Neigh:  make([][]int, n),
	}

	// Списки смежности + сортировка соседей по минимальной метке
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if matrix[i][j] != 0 {
				g.Neigh[i] = append(g.Neigh[i], j)
			}
		}
		sort.Slice(g.Neigh[i], func(a, b int) bool {
			la := strings.ToUpper(g.Labels[g.Neigh[i][a]])
			lb := strings.ToUpper(g.Labels[g.Neigh[i][b]])
			return la < lb
		})
	}

	return g, nil
}

func (g *Graph) VertexIndex(label string) (int, bool) {
	key := strings.ToUpper(strings.TrimSpace(label))
	i, ok := g.Index[key]
	return i, ok
}

func (g *Graph) DFS(start int) []int {
	n := len(g.Labels)
	visited := make([]bool, n)
	order := make([]int, 0, n)

	type frame struct {
		v   int
		pos int
	}

	visited[start] = true
	order = append(order, start)
	stack := []frame{{v: start, pos: 0}}

	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.pos >= len(g.Neigh[top.v]) {
			stack = stack[:len(stack)-1]
			continue
		}
		to := g.Neigh[top.v][top.pos]
		top.pos++
		if !visited[to] {
			visited[to] = true
			order = append(order, to)
			stack = append(stack, frame{v: to, pos: 0})
		}
	}
	return order
}

func (g *Graph) BFS(start int) []int {
	n := len(g.Labels)
	visited := make([]bool, n)
	order := make([]int, 0, n)

	q := make([]int, 0, n)
	push := func(v int) { q = append(q, v) }
	pop := func() int {
		v := q[0]
		q = q[1:]
		return v
	}

	visited[start] = true
	push(start)

	for len(q) > 0 {
		v := pop()
		order = append(order, v)
		for _, to := range g.Neigh[v] {
			if !visited[to] {
				visited[to] = true
				push(to)
			}
		}
	}
	return order
}

func (g *Graph) FormatOrder(ord []int) string {
	out := make([]string, len(ord))
	for i, v := range ord {
		out[i] = g.Labels[v]
	}
	return strings.Join(out, " ")
}
