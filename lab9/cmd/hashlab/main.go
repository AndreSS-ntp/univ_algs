package main

import (
	"fmt"
	"time"

	"github.com/AndreSS-ntp/univ_algs/lab9/internal/cli"
	"github.com/AndreSS-ntp/univ_algs/lab9/internal/csvout"
	"github.com/AndreSS-ntp/univ_algs/lab9/internal/generator"
	"github.com/AndreSS-ntp/univ_algs/lab9/internal/prehash"
	"github.com/AndreSS-ntp/univ_algs/lab9/internal/sim"
)

func main() {
	fmt.Println("ЛР: Оценка качества хеш-функций")
	fmt.Println("Сгенерируем N случайных ключей длиной 6 и оценим коллизии для 4 комбинаций методов.")
	fmt.Println()

	mSize := cli.ReadInt("Введите размер таблицы M (>0): ", 1, 1_000_000)
	nKeys := cli.ReadInt("Введите количество ключей N (>0): ", 1, 10_000_000)

	sep := cli.ReadRune("Введите разделитель CSV (Enter = ';'): ", ';')
	outPath := cli.ReadString("Введите имя выходного CSV файла (Enter = results.csv): ", "results.csv")

	useCustomSeed := cli.ReadYesNo("Задать seed вручную? (y/n, Enter = n): ", false)
	var seed int64
	if useCustomSeed {
		seed = cli.ReadInt64("Введите seed (целое число): ")
	} else {
		seed = time.Now().UnixNano()
	}

	fmt.Println()
	fmt.Printf("Параметры: M=%d, N=%d, sep='%c', out='%s', seed=%d\n", mSize, nKeys, sep, outPath, seed)
	fmt.Println("Выполняю моделирование...")

	keys := generator.GenerateKeys(nKeys, seed)
	r := prehash.MakeXorSalt(6, seed+1337)

	res := sim.Run(keys, mSize, r)

	if err := csvout.WriteResults(outPath, sep, res); err != nil {
		fmt.Println("Ошибка записи CSV:", err)
		return
	}

	fmt.Println()
	fmt.Println("Готово. Итого коллизий:")
	fmt.Printf("  Деление + Аддитивный:   %d\n", res.Total.DivAdd)
	fmt.Printf("  Деление + XOR:          %d\n", res.Total.DivXor)
	fmt.Printf("  Умножение + Аддитивный: %d\n", res.Total.MulAdd)
	fmt.Printf("  Умножение + XOR:        %d\n", res.Total.MulXor)
	fmt.Println()
	fmt.Println("Файл сохранён:", outPath)
}
