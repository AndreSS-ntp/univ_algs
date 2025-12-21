package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var reader = bufio.NewReader(os.Stdin)

func ReadString(prompt string, def string) string {
	for {
		fmt.Print(prompt)
		line, err := reader.ReadString('\n')
		if err != nil {
			// если stdin закрылся — вернём дефолт
			return def
		}
		line = strings.TrimSpace(line)
		if line == "" {
			return def
		}
		return line
	}
}

func ReadInt(prompt string, min int, max int) int {
	for {
		fmt.Print(prompt)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		v, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("  Ошибка: введите целое число.")
			continue
		}
		if v < min || v > max {
			fmt.Printf("  Ошибка: число должно быть в диапазоне [%d..%d]\n", min, max)
			continue
		}
		return v
	}
}

func ReadInt64(prompt string) int64 {
	for {
		fmt.Print(prompt)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		v, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			fmt.Println("  Ошибка: введите целое число (int64).")
			continue
		}
		return v
	}
}

func ReadYesNo(prompt string, def bool) bool {
	for {
		fmt.Print(prompt)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(strings.ToLower(line))

		if line == "" {
			return def
		}
		if line == "y" || line == "yes" || line == "д" || line == "да" {
			return true
		}
		if line == "n" || line == "no" || line == "н" || line == "нет" {
			return false
		}
		fmt.Println("  Введите y/n (да/нет).")
	}
}

func ReadRune(prompt string, def rune) rune {
	for {
		fmt.Print(prompt)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			return def
		}
		r, size := utf8.DecodeRuneInString(line)
		if r == utf8.RuneError || size == 0 {
			fmt.Println("  Ошибка: введите один символ.")
			continue
		}
		return r
	}
}
