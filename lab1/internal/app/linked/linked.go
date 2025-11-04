package linked

import (
	"fmt"

	"github.com/AndreSS-ntp/univ_algs/lab1/internal/domain"
)

//
// ===== Реализация №2: Очередь на односвязном списке (связанная память) =====
//

type node struct {
	val  domain.ElType
	next *node
}

type LinkedQueue struct {
	head, tail *node
	lastErr    error
}

var _ domain.Queue = (*LinkedQueue)(nil)

func (q *LinkedQueue) Init() {
	q.head, q.tail = nil, nil
	q.lastErr = nil
}
func (q *LinkedQueue) Empty() bool { return q.head == nil }
func (q *LinkedQueue) Full() bool  { return false }

func (q *LinkedQueue) Enqueue(x domain.ElType) bool {
	q.lastErr = nil
	n, err := allocateNode(x)
	if err != nil {
		q.lastErr = err
		return false
	}
	if q.tail == nil {
		q.head, q.tail = n, n
	} else {
		q.tail.next = n
		q.tail = n
	}
	return true
}

// LastError сообщает о последней ошибке, возникшей при попытке добавления.
// В случае нехватки памяти возвращается ненулевая ошибка, иначе nil.
func (q *LinkedQueue) LastError() error {
	return q.lastErr
}

func (q *LinkedQueue) Dequeue() (domain.ElType, bool) {
	if q.head == nil {
		return nil, false
	}
	n := q.head
	q.head = n.next
	if q.head == nil {
		q.tail = nil
	}
	return n.val, true
}

func (q *LinkedQueue) Front() (domain.ElType, bool) {
	if q.head == nil {
		return nil, false
	}
	return q.head.val, true
}

func (q *LinkedQueue) Items() []domain.ElType {
	res := make([]domain.ElType, 0)
	for n := q.head; n != nil; n = n.next {
		res = append(res, n.val)
	}
	return res
}

// FillUntilMemoryExhausted циклично добавляет элементы в очередь, пока не
// закончится память. Возвращает количество успешно добавленных элементов и
// ошибку, из-за которой цикл остановился. В реальной программе такую функцию
// не стоит вызывать без необходимости — она предназначена для имитации
// переполнения памяти по требованию преподавателя.
func (q *LinkedQueue) FillUntilMemoryExhausted(generator func(int) domain.ElType) (int, error) {
	if generator == nil {
		generator = func(i int) domain.ElType {
			return domain.NewPart(fmt.Sprintf("A%03d", i%1000), 1)
		}
	}
	count := 0
	for {
		item := generator(count)
		if !q.Enqueue(item) {
			return count, q.lastErr
		}
		count++
	}
}

func allocateNode(x domain.ElType) (_ *node, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("не удалось выделить память под элемент очереди: %v", r)
		}
	}()
	return &node{val: x}, nil
}
