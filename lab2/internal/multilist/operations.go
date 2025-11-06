package multilist

import (
	"fmt"
	"strings"
)

func New(universityCity string) *MultiList {
	return &MultiList{UniversityCity: universityCity}
}

func (ml *MultiList) AddApplicant(a *Node) {

	ml.insertAtEnd(&ml.All, a, func(n, next *Node) { n.NextAll = next })

	if a.Data.Exam1 == 5 && a.Data.Exam2 == 5 && a.Data.Exam3 == 5 {
		ml.insertAtEnd(&ml.Excellent, a, func(n, next *Node) { n.NextExcellent = next })
	}

	if a.Data.Distinct {
		ml.insertAtEnd(&ml.Distinct, a, func(n, next *Node) { n.NextDistinct = next })
	}

	if !strings.EqualFold(strings.TrimSpace(a.Data.City), strings.TrimSpace(ml.UniversityCity)) {
		ml.insertAtEnd(&ml.OutOfTown, a, func(n, next *Node) { n.NextOutOfTown = next })
	}

	if a.Data.NeedsDorm {
		ml.insertAtEnd(&ml.NeedsDorm, a, func(n, next *Node) { n.NextNeedsDorm = next })
	}
}

func (ml *MultiList) insertAtEnd(ld *ListDescriptor, node *Node, setNext func(*Node, *Node)) {
	if ld.First == nil {
		ld.First = node
		ld.Last = node
		setNext(node, nil)
	} else {
		setNext(ld.Last, node)
		setNext(node, nil)
		ld.Last = node
	}
}

func (ml *MultiList) PrintAll() {
	fmt.Println("=== Полный список абитуриентов ===")
	if ml.All.First == nil {
		fmt.Println("Список пуст.")
		return
	}
	cur := ml.All.First
	for cur != nil {
		fmt.Printf("%s | %d %d %d | Отличие: %v | %s | Общежитие: %v\n",
			cur.Data.LastName,
			cur.Data.Exam1, cur.Data.Exam2, cur.Data.Exam3,
			cur.Data.Distinct, cur.Data.City, cur.Data.NeedsDorm)
		cur = cur.NextAll
	}
}

func (ml *MultiList) PrintByDescriptor(ld *ListDescriptor, getNext func(*Node) *Node) {
	if ld.First == nil {
		fmt.Println("(Список пуст)")
		return
	}
	cur := ld.First
	for cur != nil {
		fmt.Printf("%s | %d %d %d | Отличие: %v | %s | Общежитие: %v\n",
			cur.Data.LastName,
			cur.Data.Exam1, cur.Data.Exam2, cur.Data.Exam3,
			cur.Data.Distinct, cur.Data.City, cur.Data.NeedsDorm)
		cur = getNext(cur)
	}
}

//strings.ToLower(strings.TrimSpace(a.Data.City)) != strings.ToLower(strings.TrimSpace(ml.UniversityCity))
