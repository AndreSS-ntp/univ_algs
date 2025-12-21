package avl

func Search(root *Node, key int) int {
	steps := 0
	current := root

	for current != nil {
		steps++

		if key == current.Key {
			return steps
		} else if key < current.Key {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	return steps
}
