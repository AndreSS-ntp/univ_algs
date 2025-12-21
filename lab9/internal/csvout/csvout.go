package csvout

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/AndreSS-ntp/univ_algs/lab9/internal/sim"
)

// WriteResults пишет CSV для построения графиков "коллизии vs адрес".
// Последней строкой добавляет "ИТОГО".
func WriteResults(path string, sep rune, res sim.Results) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = sep

	header := []string{
		"Адрес",
		"Коллизии_Деление_Аддитивный",
		"Коллизии_Деление_XOR",
		"Коллизии_Умножение_Аддитивный",
		"Коллизии_Умножение_XOR",
	}
	if err := w.Write(header); err != nil {
		return err
	}

	m := len(res.DivAdd)
	for i := 0; i < m; i++ {
		row := []string{
			strconv.Itoa(i),
			strconv.FormatUint(res.DivAdd[i], 10),
			strconv.FormatUint(res.DivXor[i], 10),
			strconv.FormatUint(res.MulAdd[i], 10),
			strconv.FormatUint(res.MulXor[i], 10),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	totalRow := []string{
		"ИТОГО",
		fmt.Sprintf("%d", res.Total.DivAdd),
		fmt.Sprintf("%d", res.Total.DivXor),
		fmt.Sprintf("%d", res.Total.MulAdd),
		fmt.Sprintf("%d", res.Total.MulXor),
	}
	if err := w.Write(totalRow); err != nil {
		return err
	}

	w.Flush()
	return w.Error()
}
