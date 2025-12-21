package prehash

import "math/rand"

// Additive: сумма ASCII-кодов символов строки.
func Additive(s string) uint64 {
	var sum uint64
	for i := 0; i < len(s); i++ {
		sum += uint64(s[i])
	}
	return sum
}

// XOR method (как в методичке):
// sum <- sum + ord(s[i]) xor r[i]
// r[i] фиксируется на весь эксперимент и НЕ меняется.
func XorSum(s string, r []byte) uint64 {
	var sum uint64
	limit := len(s)
	if len(r) < limit {
		limit = len(r)
	}
	for i := 0; i < limit; i++ {
		sum += uint64(s[i] ^ r[i])
	}
	// Если строка длиннее r (в этой лабе не нужно, но на всякий случай)
	for i := limit; i < len(s); i++ {
		sum += uint64(s[i])
	}
	return sum
}

// MakeXorSalt создаёт массив r длины n (обычно 6) из случайных байт.
func MakeXorSalt(n int, seed int64) []byte {
	rng := rand.New(rand.NewSource(seed))
	r := make([]byte, n)
	for i := range r {
		r[i] = byte(rng.Intn(256))
	}
	return r
}
