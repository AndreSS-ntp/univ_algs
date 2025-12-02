package pkg

import "github.com/AndreSS-ntp/univ_algs/lab5/internal/domain"

// ComputeCost вычисляет цену поиска для заданного дерева.
func ComputeCost(root *domain.Node, keys []int, p []int, q []int) int {
	n := len(keys)
	total := 0

	// Успешный поиск по всем ключам
	for i := 1; i <= n; i++ {
		key := keys[i-1]
		comparisons := searchComparisons(root, key)
		total += p[i] * comparisons
	}

	// Неуспешный поиск:
	// q[0] - значение меньше первого ключа
	if n == 0 {
		return 0
	}

	val := keys[0] - 1
	comparisons := searchComparisons(root, val)
	total += q[0] * comparisons

	// q[j] для значений между key[j] и key[j+1]
	for j := 1; j < n; j++ {
		val = (keys[j-1] + keys[j]) / 2
		comparisons = searchComparisons(root, val)
		total += q[j] * comparisons
	}

	// q[N] - значение больше последнего ключа
	val = keys[n-1] + 1
	comparisons = searchComparisons(root, val)
	total += q[n] * comparisons

	return total
}

// searchComparisons возвращает количество сравнений ключа key с узлами дерева
// при стандартном двоичном поиске.
func searchComparisons(root *domain.Node, key int) int {
	comparisons := 0
	cur := root

	for cur != nil {
		comparisons++

		if key == cur.Key {
			return comparisons
		} else if key < cur.Key {
			cur = cur.Left
		} else {
			cur = cur.Right
		}
	}

	// дошли до пустого указателя — неуспешный поиск
	return comparisons
}
