package linked

import (
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/domain"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/pkg"
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
}

const memoryReserve = 500 * 1024 * 1024

var _ domain.Queue = (*LinkedQueue)(nil)

func (q *LinkedQueue) Init() {
	q.head, q.tail = nil, nil
}
func (q *LinkedQueue) Empty() bool { return q.head == nil }
func (q *LinkedQueue) Full() bool  { return false }

func (q *LinkedQueue) Enqueue(x domain.ElType) bool {
	if headroom, err := pkg.MemoryHeadroom(); err == nil {
		if headroom <= memoryReserve {
			return false
		}
	}
	n := &node{val: x}
	if q.tail == nil {
		q.head, q.tail = n, n
	} else {
		q.tail.next = n
		q.tail = n
	}
	return true
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
