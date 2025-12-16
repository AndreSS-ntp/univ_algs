package optimal

import "github.com/AndreSS-ntp/univ_algs/lab5/internal/domain"

func BuildOptimal(keys []int, p []int, q []int) (*domain.Node, int) {
	n := len(keys)

	w := make([][]int, n+1)
	c := make([][]int, n+1)
	r := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		w[i] = make([]int, n+1)
		c[i] = make([]int, n+1)
		r[i] = make([]int, n+1)
	}

	// Базовые случаи
	for i := 0; i <= n; i++ {
		w[i][i] = q[i]
		c[i][i] = 0
	}

	for length := 1; length <= n; length++ {
		for i := 0; i <= n-length; i++ {
			j := i + length

			w[i][j] = w[i][j-1] + p[j] + q[j]

			bestCost := int(1e9)
			bestK := i + 1
			for k := i + 1; k <= j; k++ {
				cost := c[i][k-1] + c[k][j]
				if cost < bestCost {
					bestCost = cost
					bestK = k
				}
			}

			c[i][j] = w[i][j] + bestCost
			r[i][j] = bestK
		}
	}

	root := buildOptimalTreeFromRoots(keys, r, 0, n)
	return root, c[0][n]
}

func buildOptimalTreeFromRoots(keys []int, r [][]int, i, j int) *domain.Node {
	if i == j {
		return nil
	}

	k := r[i][j]
	key := keys[k-1]
	root := domain.NewNode(key)

	root.Left = buildOptimalTreeFromRoots(keys, r, i, k-1)
	root.Right = buildOptimalTreeFromRoots(keys, r, k, j)

	return root
}
