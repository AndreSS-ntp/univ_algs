package linear

import (
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/config"
	"github.com/AndreSS-ntp/univ_algs/lab1/internal/domain"
)

//
// ===== Реализация №1: Кольцевая очередь на массиве (последовательная память) =====
//
// Теория: индексы head (голова) и tail (следующая свободная позиция).
// Пусто: head == tail.
// Полно: (tail + 1) % cap == head.
// Как в учебнике: один элемент массива остаётся «запасным», поэтому для
// фактической вместимости N = cap-1. Чтобы реально хранить до 5 деталей,
// берём cap = 6.
//

type ArrayQueue struct {
	buf        [config.RingCapacity]domain.ElType
	head, tail int
}

var _ domain.Queue = (*ArrayQueue)(nil)

func (q *ArrayQueue) Init() {
	q.head, q.tail = 0, 0
	for i := range q.buf {
		q.buf[i] = nil
	}
}
func (q *ArrayQueue) Empty() bool { return q.head == q.tail }
func (q *ArrayQueue) Full() bool  { return (q.tail+1)%config.RingCapacity == q.head }

func (q *ArrayQueue) Enqueue(x domain.ElType) bool {
	if q.Full() {
		return false
	}
	q.buf[q.tail] = x
	q.tail = (q.tail + 1) % config.RingCapacity
	return true
}

func (q *ArrayQueue) Dequeue() (domain.ElType, bool) {
	if q.Empty() {
		return nil, false
	}
	x := q.buf[q.head]
	q.buf[q.head] = nil
	q.head = (q.head + 1) % config.RingCapacity
	return x, true
}

func (q *ArrayQueue) Front() (domain.ElType, bool) {
	if q.Empty() {
		return nil, false
	}
	return q.buf[q.head], true
}

func (q *ArrayQueue) Items() []domain.ElType {
	res := make([]domain.ElType, 0)
	i := q.head
	for i != q.tail {
		res = append(res, q.buf[i])
		i = (i + 1) % config.RingCapacity
	}
	return res
}
