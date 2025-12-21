package experiment

import (
	"math/rand"
	"time"
)

func GenerateData(n int) []int {
	rand.Seed(time.Now().UnixNano())
	data := make([]int, n)

	for i := 0; i < n; i++ {
		data[i] = rand.Intn(n * 10)
	}
	return data
}
