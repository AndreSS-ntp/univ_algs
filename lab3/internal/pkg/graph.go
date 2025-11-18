package pkg

import "fmt"

// Graph - ориентированный граф в виде списка смежности.
type Graph struct {
	vertices []rune          // метки вершин
	adj      map[rune][]rune // список смежности: из вершины -> список вершин
}

// NewGraph создаёт пустой граф.
func NewGraph() *Graph {
	return &Graph{
		adj: make(map[rune][]rune),
	}
}

// AddVertex добавляет вершину с заданной меткой.
// Возвращает false, если вершина уже есть.
func (g *Graph) AddVertex(v rune) bool {
	if _, ok := g.adj[v]; ok {
		return false
	}
	g.vertices = append(g.vertices, v)
	g.adj[v] = []rune{}
	return true
}

// AddEdge добавляет ориентированное ребро from -> to.
func (g *Graph) AddEdge(from, to rune) error {
	if _, ok := g.adj[from]; !ok {
		return fmt.Errorf("вершина %c не существует", from)
	}
	if _, ok := g.adj[to]; !ok {
		return fmt.Errorf("вершина %c не существует", to)
	}
	// проверяем, нет ли уже такого ребра
	for _, v := range g.adj[from] {
		if v == to {
			return nil // уже есть
		}
	}
	g.adj[from] = append(g.adj[from], to)
	return nil
}

// RemoveVertex удаляет вершину и все инцидентные ей рёбра.
func (g *Graph) RemoveVertex(v rune) bool {
	if _, ok := g.adj[v]; !ok {
		return false
	}

	// убираем из списка вершин
	idx := -1
	for i, x := range g.vertices {
		if x == v {
			idx = i
			break
		}
	}
	if idx >= 0 {
		g.vertices = append(g.vertices[:idx], g.vertices[idx+1:]...)
	}

	// удаляем все исходящие рёбра
	delete(g.adj, v)

	// удаляем все входящие рёбра
	for from, list := range g.adj {
		newList := make([]rune, 0, len(list))
		for _, to := range list {
			if to != v {
				newList = append(newList, to)
			}
		}
		g.adj[from] = newList
	}

	return true
}

// HasEdge проверяет наличие ребра from -> to.
func (g *Graph) HasEdge(from, to rune) bool {
	list, ok := g.adj[from]
	if !ok {
		return false
	}
	for _, v := range list {
		if v == to {
			return true
		}
	}
	return false
}

// Print выводит граф в виде списка смежности.
func (g *Graph) Print() {
	fmt.Println("Список смежности (ориентированный граф):")
	for _, v := range g.vertices {
		fmt.Printf("%c: ", v)
		for _, to := range g.adj[v] {
			fmt.Printf("%c ", to)
		}
		fmt.Println()
	}
}
