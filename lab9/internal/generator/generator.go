package generator

import "math/rand"

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateKeys генерирует N случайных строк длиной 6 из [0-9A-Za-z].
// seed задаёт воспроизводимость.
func GenerateKeys(n int, seed int64) []string {
	rng := rand.New(rand.NewSource(seed))
	keys := make([]string, n)

	for i := 0; i < n; i++ {
		b := make([]byte, 6)
		for j := 0; j < 6; j++ {
			b[j] = charset[rng.Intn(len(charset))]
		}
		keys[i] = string(b)
	}
	return keys
}
