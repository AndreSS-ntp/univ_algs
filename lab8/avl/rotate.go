package avl

func rotateRight(y *Node) *Node {
	x := y.Left
	T2 := x.Right

	x.Right = y
	y.Left = T2

	y.Height = max(height(y.Left), height(y.Right)) + 1
	x.Height = max(height(x.Left), height(x.Right)) + 1

	return x
}

func rotateLeft(x *Node) *Node {
	y := x.Right
	T2 := y.Left

	y.Left = x
	x.Right = T2

	x.Height = max(height(x.Left), height(x.Right)) + 1
	y.Height = max(height(y.Left), height(y.Right)) + 1

	return y
}
