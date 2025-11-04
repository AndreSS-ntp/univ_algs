package linked

import (
	"errors"
	"fmt"
	"unsafe"

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
	approxMem  uint64
	memLimit   uint64
}

var _ domain.Queue = (*LinkedQueue)(nil)

func (q *LinkedQueue) Init() {
	q.head, q.tail = nil, nil
	q.lastErr = nil
	q.approxMem = 0
}
func (q *LinkedQueue) Empty() bool { return q.head == nil }
func (q *LinkedQueue) Full() bool  { return false }

func (q *LinkedQueue) Enqueue(x domain.ElType) bool {
	q.lastErr = nil
	if !q.reserveMemory() {
		return false
	}
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
	q.approxMem += perElementFootprint
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
	q.releaseMemory()
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

// SetMemoryLimit устанавливает мягкий предел памяти в байтах для очереди.
// Если limit == 0, ограничение отключено.
func (q *LinkedQueue) SetMemoryLimit(limit uint64) {
	q.memLimit = limit
}

// MemoryLimit возвращает текущий установленный предел памяти (0 — без ограничений).
func (q *LinkedQueue) MemoryLimit() uint64 {
	return q.memLimit
}

// ApproxMemoryUsage возвращает приблизительный объём памяти, занятый элементами.
func (q *LinkedQueue) ApproxMemoryUsage() uint64 {
	return q.approxMem
}

var errMemoryLimit = errors.New("недостаточно памяти: достигнут установленный предел")

// DefaultStressMemoryLimit — предельный объём памяти (байт), который используется
// стресс-тестом, чтобы избежать аварийного завершения программы на машинах
// с большим объёмом ОЗУ.
const DefaultStressMemoryLimit = 32 << 20

const (
	nodeFootprint       = uint64(unsafe.Sizeof(node{}))
	partFootprint       = uint64(unsafe.Sizeof(domain.Part{}))
	perElementFootprint = nodeFootprint + partFootprint
)

func (q *LinkedQueue) reserveMemory() bool {
	if q.memLimit == 0 {
		return true
	}
	if q.approxMem+perElementFootprint > q.memLimit {
		q.lastErr = errMemoryLimit
		return false
	}
	return true
}

func (q *LinkedQueue) releaseMemory() {
	if q.approxMem >= perElementFootprint {
		q.approxMem -= perElementFootprint
	} else {
		q.approxMem = 0
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
