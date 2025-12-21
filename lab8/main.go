package main

import (
	"avl-lab/experiment"
	"fmt"
)

func main() {
	sizes := []int{100, 1000, 10000, 20000, 50000, 100000}

	fmt.Println("N\tAverage search path")
	for _, n := range sizes {
		data := experiment.GenerateData(n)
		avg := experiment.AverageSearchPath(data)
		fmt.Printf("%d\t%.2f\n", n, avg)
	}
}
