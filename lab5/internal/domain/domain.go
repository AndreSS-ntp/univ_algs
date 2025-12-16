package domain

import (
	"fmt"
	"strings"
)

type Node struct {
	Key         int
	Left, Right *Node
}

func NewNode(key int) *Node {
	return &Node{Key: key}
}

func PrintTree(root *Node, level int) {
	if root == nil {
		return
	}

	PrintTree(root.Right, level+1)

	fmt.Printf("%s%d\n", strings.Repeat("  ", level), root.Key)

	PrintTree(root.Left, level+1)
}
