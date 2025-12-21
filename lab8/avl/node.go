package avl

type Node struct {
	Key    int
	Height int
	Left   *Node
	Right  *Node
}

func height(n *Node) int {
	if n == nil {
		return 0
	}
	return n.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getBalance(n *Node) int {
	if n == nil {
		return 0
	}
	return height(n.Right) - height(n.Left)
}
