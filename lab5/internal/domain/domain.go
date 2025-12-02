package domain

import (
	"fmt"
	"strings"
)

// Узел бинарного дерева поиска
type Node struct {
	Key         int
	Left, Right *Node
}

func NewNode(key int) *Node {
	return &Node{Key: key}
}

// Печать дерева "справа налево" (правое поддерево наверху, левое внизу)
func PrintTree(root *Node, level int) {
	if root == nil {
		return
	}

	// сначала правое поддерево
	PrintTree(root.Right, level+1)

	// сам узел
	fmt.Printf("%s%d\n", strings.Repeat("  ", level), root.Key)

	// потом левое поддерево
	PrintTree(root.Left, level+1)
}
