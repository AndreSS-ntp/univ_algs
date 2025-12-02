package balanced

import "github.com/AndreSS-ntp/univ_algs/lab5/internal/domain"

// Построение полностью сбалансированного дерева из отсортированного массива ключей.
func BuildBalanced(keys []int) *domain.Node {
	return buildBalancedRec(keys, 0, len(keys)-1)
}

func buildBalancedRec(keys []int, left, right int) *domain.Node {
	if left > right {
		return nil
	}

	mid := (left + right) / 2
	root := domain.NewNode(keys[mid])

	root.Left = buildBalancedRec(keys, left, mid-1)
	root.Right = buildBalancedRec(keys, mid+1, right)

	return root
}
