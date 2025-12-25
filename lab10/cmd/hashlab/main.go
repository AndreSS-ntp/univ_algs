package main

import (
	"fmt"
	"time"

	"github.com/AndreSS-ntp/univ_algs/lab10/internal/cli"
	"github.com/AndreSS-ntp/univ_algs/lab10/internal/csvout10"
	"github.com/AndreSS-ntp/univ_algs/lab10/internal/lab10"
	"github.com/AndreSS-ntp/univ_algs/lab10/internal/prehash"
)

func main() {
	fmt.Println("ЛР №10: Анализ методов разрешения коллизий (открытая адресация)")
	fmt.Println("Хеширование: деление/умножение. Опробование: линейное/квадратичное.")
	fmt.Println("Преобразование ключей: XOR (обязательно).")
	fmt.Println()

	M := cli.ReadInt("Введите размер таблицы M (>=10): ", 10, 1_000_000)
	E := cli.ReadInt("Введите число экспериментов E (например 10..100): ", 1, 10_000)

	sep := cli.ReadRune("Введите разделитель CSV (Enter = ';'): ", ';')
	outPath := cli.ReadString("Введите имя выходного CSV файла (Enter = lab10_results.csv): ", "lab10_results.csv")

	useCustomSeed := cli.ReadYesNo("Задать seed вручную? (y/n, Enter = n): ", false)
	var seed int64
	if useCustomSeed {
		seed = cli.ReadInt64("Введите seed (целое число): ")
	} else {
		seed = time.Now().UnixNano()
	}

	// Константы опробования (можно оставить по умолчанию).
	// Для линейного: h_i = (h0 + c*i) mod M
	// Для квадратичного: h_i = (h0 + c*i + d*i*i) mod M
	defC, defD := lab10.DefaultProbeConstants(M)
	fmt.Println()
	fmt.Printf("Рекомендуемые константы: c=%d, d=%d\n", defC, defD)

	c := cli.ReadInt(fmt.Sprintf("Введите c (Enter = %d): ", defC), 1, 1_000_000_000)
	d := cli.ReadInt(fmt.Sprintf("Введите d (Enter = %d): ", defD), 0, 1_000_000_000)

	// Исправление/проверка c для линейного (желательно gcd(c, M)=1)
	cFixed := lab10.AdjustCToCoprime(c, M)
	if cFixed != c {
		fmt.Printf("Внимание: c=%d не взаимно просто с M=%d. Использую ближайшее c=%d.\n", c, M, cFixed)
		c = cFixed
	}

	// r для XOR должен быть фиксирован (чтобы эксперимент был корректнее)
	r := prehash.MakeXorSalt(6, seed+1337)

	fmt.Println()
	fmt.Printf("Параметры: M=%d, E=%d, c=%d, d=%d, sep='%c', out='%s', seed=%d\n", M, E, c, d, sep, outPath, seed)
	fmt.Println("Запуск экспериментов...")

	rows := lab10.RunAll(M, E, c, d, seed, r)

	if err := csvout10.Write(outPath, sep, rows); err != nil {
		fmt.Println("Ошибка записи CSV:", err)
		return
	}

	fmt.Println("Готово. CSV сохранён:", outPath)
}
