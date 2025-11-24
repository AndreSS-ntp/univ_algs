package pkg

import "fmt"

type Graph struct {
	vertices []rune
	adj      [][]rune
}

func NewGraph() *Graph {
	return &Graph{}
}

func (g *Graph) findVertexIndex(v rune) int {
	for i, label := range g.vertices {
		if label == v {
			return i
		}
	}
	return -1
}

func (g *Graph) AddVertex(v rune) bool {
	if g.findVertexIndex(v) != -1 {
		return false
	}
	g.vertices = append(g.vertices, v)
	g.adj = append(g.adj, []rune{})
	return true
}

func (g *Graph) AddEdge(from, to rune) error {
	fromIdx := g.findVertexIndex(from)
	if fromIdx == -1 {
		return fmt.Errorf("вершина %c не существует", from)
	}
	if g.findVertexIndex(to) == -1 {
		return fmt.Errorf("вершина %c не существует", to)
	}

	for _, v := range g.adj[fromIdx] {
		if v == to {
			return nil
		}
	}

	g.adj[fromIdx] = append(g.adj[fromIdx], to)
	return nil
}

func (g *Graph) RemoveVertex(v rune) bool {
	idx := g.findVertexIndex(v)
	if idx == -1 {
		return false
	}

	g.vertices = append(g.vertices[:idx], g.vertices[idx+1:]...)

	g.adj = append(g.adj[:idx], g.adj[idx+1:]...)

	for i := range g.adj {
		list := g.adj[i]
		newList := make([]rune, 0, len(list))
		for _, to := range list {
			if to != v {
				newList = append(newList, to)
			}
		}
		g.adj[i] = newList
	}

	return true
}

func (g *Graph) HasEdge(from, to rune) bool {
	fromIdx := g.findVertexIndex(from)
	if fromIdx == -1 {
		return false
	}
	for _, v := range g.adj[fromIdx] {
		if v == to {
			return true
		}
	}
	return false
}

func (g *Graph) Print() {
	fmt.Println("Список смежности (ориентированный граф):")
	for i, v := range g.vertices {
		fmt.Printf("%c: ", v)
		for _, to := range g.adj[i] {
			fmt.Printf("%c ", to)
		}
		fmt.Println()
	}
}
