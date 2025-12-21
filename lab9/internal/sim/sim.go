package sim

import (
	"github.com/AndreSS-ntp/univ_algs/lab9/internal/hasher"
	"github.com/AndreSS-ntp/univ_algs/lab9/internal/prehash"
)

type Totals struct {
	DivAdd uint64
	DivXor uint64
	MulAdd uint64
	MulXor uint64
}

type Results struct {
	// Коллизии по адресам:
	DivAdd []uint64
	DivXor []uint64
	MulAdd []uint64
	MulXor []uint64

	Total Totals
}

// Run моделирует рассеивание ключей и считает коллизии на каждый адрес
// для 4 комбинаций: (деление/умножение) x (аддитивный/XOR).
func Run(keys []string, m int, r []byte) Results {
	// Счётчики количества ключей по адресам (потом переводим в коллизии)
	divAddCnt := make([]uint64, m)
	divXorCnt := make([]uint64, m)
	mulAddCnt := make([]uint64, m)
	mulXorCnt := make([]uint64, m)

	for _, s := range keys {
		kAdd := prehash.Additive(s)
		kXor := prehash.XorSum(s, r)

		a1 := hasher.DivisionHash(kAdd, m)
		a2 := hasher.DivisionHash(kXor, m)
		a3 := hasher.MultiplicationHash(kAdd, m)
		a4 := hasher.MultiplicationHash(kXor, m)

		divAddCnt[a1]++
		divXorCnt[a2]++
		mulAddCnt[a3]++
		mulXorCnt[a4]++
	}

	// Переводим counts -> collisions: max(0, count-1)
	res := Results{
		DivAdd: make([]uint64, m),
		DivXor: make([]uint64, m),
		MulAdd: make([]uint64, m),
		MulXor: make([]uint64, m),
	}

	for i := 0; i < m; i++ {
		res.DivAdd[i] = toCollisions(divAddCnt[i])
		res.DivXor[i] = toCollisions(divXorCnt[i])
		res.MulAdd[i] = toCollisions(mulAddCnt[i])
		res.MulXor[i] = toCollisions(mulXorCnt[i])

		res.Total.DivAdd += res.DivAdd[i]
		res.Total.DivXor += res.DivXor[i]
		res.Total.MulAdd += res.MulAdd[i]
		res.Total.MulXor += res.MulXor[i]
	}

	return res
}

func toCollisions(cnt uint64) uint64 {
	if cnt == 0 {
		return 0
	}
	if cnt == 1 {
		return 0
	}
	return cnt - 1
}
