package experiment

import "avl-lab/avl"

func AverageSearchPath(data []int) float64 {
	var root *avl.Node

	for _, v := range data {
		root = avl.Insert(root, v)
	}

	totalSteps := 0
	for _, v := range data {
		totalSteps += avl.Search(root, v)
	}

	return float64(totalSteps) / float64(len(data))
}
