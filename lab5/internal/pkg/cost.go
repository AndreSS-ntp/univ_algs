package pkg

import "github.com/AndreSS-ntp/univ_algs/lab5/internal/domain"

func ComputeCost(root *domain.Node, keys []int, p []int, q []int) int {
	n := len(keys)
	total := 0

	for i := 1; i <= n; i++ {
		key := keys[i-1]
		comparisons := searchComparisons(root, key)
		total += p[i] * comparisons
	}

	if n == 0 {
		return 0
	}

	val := keys[0] - 1
	comparisons := searchComparisons(root, val)
	total += q[0] * comparisons

	for j := 1; j < n; j++ {
		val = (keys[j-1] + keys[j]) / 2
		comparisons = searchComparisons(root, val)
		total += q[j] * comparisons
	}

	val = keys[n-1] + 1
	comparisons = searchComparisons(root, val)
	total += q[n] * comparisons

	return total
}

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

	return comparisons
}
