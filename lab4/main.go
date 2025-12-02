package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lab/internal/ioutil"
	"lab/model"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	readNum := func() (int, error) {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		num, err := strconv.Atoi(line)

		return num, err
	}

	var root *model.BinaryTree
	for {
		fmt.Println("\nМеню:")
		fmt.Println("1 - Построить дерево по данным из файла")
		fmt.Println("2 - Вывести дерево")
		fmt.Println("3 - Поиск вершин, где разница между потомками поддеревьей равна единице")
		fmt.Println("4 - Высота дерева (итеративная)")
		fmt.Println("5 - Найти k-й лист слева направо")
		fmt.Println("6 - Удалить узел по ключу")
		fmt.Println("7 - Добавить ключ")
		fmt.Println("0 - Выход")
		fmt.Print("Выбор: ")

		choiceLine, _ := reader.ReadString('\n')
		choiceLine = strings.TrimSpace(choiceLine)
		choice, _ := strconv.Atoi(choiceLine)

		switch choice {
		case 1:
			fmt.Print("Введите путь к входному файлу с целыми числами: ")
			path, _ := reader.ReadString('\n')
			path = strings.TrimSpace(path)

			nums, err := ioutil.ReadIntsFromFile(path)
			if err != nil {
				fmt.Printf("Ошибка чтения файла: %v\n", err)
				return
			}

			for _, v := range nums {
				root = model.InsertNode(root, v)
			}
		case 2:
			fmt.Print("Бинарное дерево: ")
			model.PreorderPrint(root)
		case 3:
			found := model.FindUnbalancedNodes(root)
			fmt.Printf("Узлы, где |L-R|==1: %v\n", found)
		case 4:
			h := model.HeightIterativeDFS(root)
			fmt.Printf("Высота дерева: %d\n", h)
		case 5:
			fmt.Print("Введите k: ")
			k, err := readNum()
			if err != nil || k <= 0 {
				fmt.Println("Неверное k")
				continue
			}
			node, ok := model.KthLeafInorder(root, k)
			if ok {
				fmt.Printf("%d-й лист: %d\n", k, node.Key)
			} else {
				fmt.Println("Лист с таким номером не найден")
			}
		case 6:
			fmt.Print("Введите ключ для удаления: ")
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			key, err := strconv.Atoi(line)
			if err != nil {
				fmt.Println("Неверный ключ")
				continue
			}
			root = model.DeleteNode(root, key)
			fmt.Println("Удаление выполнено (если такой ключ был).")
		case 7:
			fmt.Print("Введите ключ для добавления: ")
			key, err := readNum()
			if err != nil {
				fmt.Println("Неверный ключ")
				continue
			}
			root = model.InsertNode(root, key)
			fmt.Println("Добавление выполнено.")
		case 0:
			fmt.Println("Выход.")
			return
		default:
			fmt.Println("Неверный выбор")
		}
	}
}
