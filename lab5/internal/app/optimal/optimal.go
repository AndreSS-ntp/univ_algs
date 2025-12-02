package optimal

import "github.com/AndreSS-ntp/univ_algs/lab5/internal/domain"

// BuildOptimal строит оптимальное дерево поиска по массивам keys, p, q.
// Возвращает корень дерева и минимальную цену поиска C(1, N).
func BuildOptimal(keys []int, p []int, q []int) (*domain.Node, int) {
	n := len(keys) // N

	// Матрицы w, c, r размером (n+1) x (n+1), индексация i,j: 0..N
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
		w[i][i] = q[i] // только один абстрактный узел
		c[i][i] = 0
	}

	// length = j - i, количество ключей в поддереве
	for length := 1; length <= n; length++ {
		for i := 0; i <= n-length; i++ {
			j := i + length

			// считаем W(i, j)
			w[i][j] = w[i][j-1] + p[j] + q[j]

			// ищем минимум по k
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

// Реконструкция дерева по матрице R(i,j)
func buildOptimalTreeFromRoots(keys []int, r [][]int, i, j int) *domain.Node {
	if i == j {
		return nil // пустое поддерево
	}

	k := r[i][j]     // индекс корня в терминах p/q (1..N)
	key := keys[k-1] // keys индексация с 0
	root := domain.NewNode(key)

	root.Left = buildOptimalTreeFromRoots(keys, r, i, k-1)
	root.Right = buildOptimalTreeFromRoots(keys, r, k, j)

	return root
}
