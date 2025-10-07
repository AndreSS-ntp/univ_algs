package pkg

import (
	"bufio"
	"fmt"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/domain"
	"strconv"
	"strings"
)

//
// ===== Вспомогательные функции =====
//

func ShowQueue(q domain.Queue, t int) {
	items := q.Items()
	if len(items) == 0 {
		fmt.Println("Очередь пуста.")
		return
	}
	fmt.Printf("Очередь на t=%d (от головы к хвосту):\n", t)
	for i, p := range items {
		fmt.Printf("%2d) %s  осталось: %d  (заявлено: %d)\n", i+1, p.Code, p.Time, p.StartTime)
	}
}

func ReadLine(in *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	text, _ := in.ReadString('\n')
	return strings.TrimSpace(text)
}

func ReadInt(in *bufio.Reader, prompt string) int {
	for {
		s := ReadLine(in, prompt)
		if s == "" {
			continue
		}
		if v, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			return v
		}
		fmt.Println("Введите целое число.")
	}
}

func ReadPositiveInt(in *bufio.Reader, prompt string) int {
	for {
		v := ReadInt(in, prompt)
		if v > 0 {
			return v
		}
		fmt.Println("Число должно быть > 0.")
	}
}

func NormalizeCode(s string) string {
	s = strings.TrimSpace(s)
	// Упростим: обрежем/дополняем до 4 символов
	r := []rune(s)
	if len(r) >= 4 {
		return string(r[:4])
	}
	for len(r) < 4 {
		r = append(r, '_')
	}
	return string(r)
}
