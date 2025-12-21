package avl

func Insert(node *Node, key int) *Node {
	if node == nil {
		return &Node{Key: key, Height: 1}
	}

	if key < node.Key {
		node.Left = Insert(node.Left, key)
	} else if key > node.Key {
		node.Right = Insert(node.Right, key)
	} else {
		return node
	}

	node.Height = 1 + max(height(node.Left), height(node.Right))
	balance := getBalance(node)

	// LL
	if balance < -1 && key < node.Left.Key {
		return rotateRight(node)
	}

	// RR
	if balance > 1 && key > node.Right.Key {
		return rotateLeft(node)
	}

	// LR
	if balance < -1 && key > node.Left.Key {
		node.Left = rotateLeft(node.Left)
		return rotateRight(node)
	}

	// RL
	if balance > 1 && key < node.Right.Key {
		node.Right = rotateRight(node.Right)
		return rotateLeft(node)
	}

	return node
}
