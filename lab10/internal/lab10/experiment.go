package lab10

import (
	"github.com/AndreSS-ntp/univ_algs/lab10/internal/generator"
	"github.com/AndreSS-ntp/univ_algs/lab10/internal/hasher"
	"github.com/AndreSS-ntp/univ_algs/lab10/internal/oa"
	"github.com/AndreSS-ntp/univ_algs/lab10/internal/prehash"
)

type Row struct {
	N       int
	DivLin  float64
	DivQuad float64
	MulLin  float64
	MulQuad float64
}

// DefaultProbeConstants — “разумные” дефолты из рекомендаций методички:
// c — достаточно большое (≈ M/5), d — небольшое (1).
func DefaultProbeConstants(M int) (c int, d int) {
	c = M / 5
	if c < 1 {
		c = 1
	}
	d = 1
	// c лучше сделать взаимно простым с M (для линейного)
	c = AdjustCToCoprime(c, M)
	return c, d
}

func AdjustCToCoprime(c int, M int) int {
	if gcd(c, M) == 1 {
		return c
	}
	// подберём вверх ближайшее c, взаимно простое с M
	for x := c + 1; x < c+M+10; x++ {
		if gcd(x, M) == 1 {
			return x
		}
	}
	// fallback
	return 1
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// RunAll строит точки для графиков: N=10..M шаг 10.
// Для каждой точки усредняет по E экспериментам.
func RunAll(M int, E int, c int, d int, seed int64, r []byte) []Row {
	pp := oa.ProbeParams{M: M, C: c, D: d}

	rows := make([]Row, 0, M/10)

	for n := 10; n <= M; n += 10 {
		var sumDivLin, sumDivQuad, sumMulLin, sumMulQuad float64

		for e := 0; e < E; e++ {
			// разные ключи на каждый эксперимент
			keys := generator.GenerateKeys(n, seed+int64(n)*10_000+int64(e))

			sumDivLin += avgProbes(keys, M, pp, oa.Linear, true, r)
			sumDivQuad += avgProbes(keys, M, pp, oa.Quadratic, true, r)

			sumMulLin += avgProbes(keys, M, pp, oa.Linear, false, r)
			sumMulQuad += avgProbes(keys, M, pp, oa.Quadratic, false, r)
		}

		row := Row{
			N:       n,
			DivLin:  sumDivLin / float64(E),
			DivQuad: sumDivQuad / float64(E),
			MulLin:  sumMulLin / float64(E),
			MulQuad: sumMulQuad / float64(E),
		}
		rows = append(rows, row)
	}

	return rows
}

// avgProbes считает среднюю длину пути поиска при вставке по всем ключам в одном эксперименте.
// useDivision=true => метод деления, false => метод умножения.
func avgProbes(keys []string, M int, pp oa.ProbeParams, kind oa.ProbeKind, useDivision bool, r []byte) float64 {
	t := oa.NewTable(M)

	var probesSum int
	var inserted int

	for _, key := range keys {
		k := prehash.XorSum(key, r)

		var h0 int
		if useDivision {
			h0 = hasher.DivisionHash(k, M)
		} else {
			h0 = hasher.MultiplicationHash(k, M)
		}

		probes, _ := t.Insert(key, h0, kind, pp)
		probesSum += probes
		inserted++
	}

	if inserted == 0 {
		return 0
	}
	return float64(probesSum) / float64(inserted)
}
