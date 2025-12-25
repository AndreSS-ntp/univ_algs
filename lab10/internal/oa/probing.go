package oa

type ProbeKind int

const (
	Linear ProbeKind = iota
	Quadratic
)

type ProbeParams struct {
	M int
	C int
	D int
}

func ProbeAddress(h0 int, i int, kind ProbeKind, p ProbeParams) int {
	switch kind {
	case Linear:
		// (h0 + c*i) mod M
		return mod(h0+p.C*i, p.M)
	case Quadratic:
		// (h0 + c*i + d*i^2) mod M
		return mod(h0+p.C*i+p.D*i*i, p.M)
	default:
		return mod(h0+i, p.M)
	}
}

func mod(x int, m int) int {
	r := x % m
	if r < 0 {
		r += m
	}
	return r
}
