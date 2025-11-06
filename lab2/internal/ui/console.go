package ui

import (
	"bufio"
	"fmt"
	"multilist/internal/model"
	"multilist/internal/multilist"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func readLine(prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func readInt(prompt string, min, max int) int {
	for {
		s := readLine(prompt)
		v, err := strconv.Atoi(s)
		if err == nil && v >= min && v <= max {
			return v
		}
		fmt.Printf("Введите число [%d..%d]\n", min, max)
	}
}

func readBool(prompt string) bool {
	for {
		s := strings.ToLower(readLine(prompt + " (y/n): "))
		switch s {
		case "y", "д", "да":
			return true
		case "n", "н", "нет":
			return false
		default:
			fmt.Println("Введите 'y' или 'n'.")
		}
	}
}

func CreateApplicantFromInput() *model.Applicant {
	fmt.Println("\n=== Добавление абитуриента ===")
	last := readLine("Фамилия: ")
	ex1 := readInt("Экзамен 1 (1-5): ", 1, 5)
	ex2 := readInt("Экзамен 2 (1-5): ", 1, 5)
	ex3 := readInt("Экзамен 3 (1-5): ", 1, 5)

	// allExcellent := ex1 == 5 && ex2 == 5 && ex3 == 5

	// var dist bool
	// if allExcellent {
	// 	dist = readBool("Аттестат с отличием?")
	// } else {
	// 	fmt.Println("Аттестат с отличием невозможен — не все экзамены на 5.")
	// 	dist = false
	// }

	dist := readBool("Аттестат с отличием?")
	city := readLine("Город проживания: ")
	dorm := readBool("Нуждается в общежитии?")
	return &model.Applicant{
		LastName:  last,
		Exam1:     ex1,
		Exam2:     ex2,
		Exam3:     ex3,
		Distinct:  dist,
		City:      city,
		NeedsDorm: dorm,
	}
}

func Run() {
	fmt.Println("=== Многосвязный список абитуриентов ===")
	uniCity := readLine("Введите город университета: ")
	ml := multilist.New(uniCity)

	for {
		fmt.Println("\n===== МЕНЮ =====")
		fmt.Println("1. Добавить абитуриента")
		fmt.Println("2. Показать полный список")
		fmt.Println("3. Показать список сдавших на отлично")
		fmt.Println("4. Показать список с аттестатом с отличием")
		fmt.Println("5. Показать список из других городов")
		fmt.Println("6. Показать список нуждающихся в общежитии")
		fmt.Println("7. Удалить абитуриента по фамилии (первое совпадение)")
		fmt.Println("8. Удалить всех абитуриентов с фамилией")
		fmt.Println("9. Удалить все записи")
		fmt.Println("0. Выход")

		switch readLine("Выбор: ") {
		case "1":
			a := CreateApplicantFromInput()
			n := &multilist.Node{Data: a}
			ml.AddApplicant(n)
			fmt.Println("Абитуриент добавлен.")

		case "2":
			ml.PrintAll()

		case "3":
			fmt.Println("=== Все экзамены на 'отлично' ===")
			ml.PrintByDescriptor(&ml.Excellent, func(n *multilist.Node) *multilist.Node { return n.NextExcellent })

		case "4":
			fmt.Println("=== Аттестат с отличием ===")
			ml.PrintByDescriptor(&ml.Distinct, func(n *multilist.Node) *multilist.Node { return n.NextDistinct })

		case "5":
			fmt.Println("=== Проживает вне города университета ===")
			ml.PrintByDescriptor(&ml.OutOfTown, func(n *multilist.Node) *multilist.Node { return n.NextOutOfTown })

		case "6":
			fmt.Println("=== Нуждается в общежитии ===")
			ml.PrintByDescriptor(&ml.NeedsDorm, func(n *multilist.Node) *multilist.Node { return n.NextNeedsDorm })

		case "7":
			last := readLine("Введите фамилию для удаления: ")
			if ml.DeleteByLastName(last) {
				fmt.Println("Абитуриент удалён.")
			} else {
				fmt.Println("Абитуриент с такой фамилией не найден.")
			}

		case "8":
			last := readLine("Введите фамилию для массового удаления: ")
			count := ml.DeleteAllByLastName(last)
			fmt.Printf("Удалено записей: %d\n", count)

		case "9":
			if readBool("Вы уверены, что хотите удалить все записи?") {
				ml.DeleteAll()
				fmt.Println("Все записи удалены.")
			}

		case "0":
			fmt.Println("Выход из программы.")
			return

		default:
			fmt.Println("Неверный пункт меню.")
		}
	}
}
