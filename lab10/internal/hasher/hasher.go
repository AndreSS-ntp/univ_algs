package hasher

import "math"

func DivisionHash(k uint64, m int) int {
	return int(k % uint64(m))
}

func MultiplicationHash(k uint64, m int) int {
	const A = 0.6180339887498949

	x := float64(k) * A
	frac := x - math.Floor(x)
	addr := int(math.Floor(float64(m) * frac))

	// страховка от редких погрешностей float
	if addr < 0 {
		return 0
	}
	if addr >= m {
		return m - 1
	}
	return addr
}
