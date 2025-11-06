package domain

type Part struct {
	Code      string // код детали (4 символа)
	Time      int    // оставшееся время обработки
	StartTime int    // исходное время (для сообщений)
}

func NewPart(code string, t int) *Part {
	return &Part{Code: code, Time: t, StartTime: t}
}

type ElType = *Part

type Queue interface {
	Init()
	Enqueue(x ElType) bool
	Dequeue() (ElType, bool)
	Empty() bool
	Full() bool            // для связной реализации всегда false
	Front() (ElType, bool) // неразрушающий доступ к первому элементу
	Items() []ElType       // содержимое очереди от головы к хвосту
}
