package csvout10

import (
	"encoding/csv"
	"fmt"
	"github.com/AndreSS-ntp/univ_algs/lab10/internal/lab10"
	"os"
)

func Write(path string, sep rune, rows []lab10.Row) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = sep

	header := []string{
		"КоличествоКлючей",
		"ДлинаПути_Деление_Линейное",
		"ДлинаПути_Деление_Квадратичное",
		"ДлинаПути_Умножение_Линейное",
		"ДлинаПути_Умножение_Квадратичное",
	}
	if err := w.Write(header); err != nil {
		return err
	}

	for _, r := range rows {
		row := []string{
			fmt.Sprintf("%d", r.N),
			formatFloatRU(r.DivLin),
			formatFloatRU(r.DivQuad),
			formatFloatRU(r.MulLin),
			formatFloatRU(r.MulQuad),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}

// Excel RU часто удобнее воспринимает дроби с запятой.
func formatFloatRU(v float64) string {
	s := fmt.Sprintf("%.4f", v)
	return s
}
