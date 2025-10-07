package domain

//
// ===== Описание элемента очереди (ElType) =====
//

type Part struct {
	Code      string // код детали (4 символа)
	Time      int    // оставшееся время обработки
	StartTime int    // исходное время (для сообщений)
}

func NewPart(code string, t int) *Part {
	return &Part{Code: code, Time: t, StartTime: t}
}

// ElType — тип элемента очереди, как просили в задании.
// Используем указатель, чтобы можно было менять оставшееся время
// у элемента на «голове» без удаления из очереди.
type ElType = *Part

//
// ===== АТД Queue (интерфейс) =====
//

type Queue interface {
	Init()
	Enqueue(x ElType) bool
	Dequeue() (ElType, bool)
	Empty() bool
	Full() bool            // для связной реализации всегда false
	Front() (ElType, bool) // неразрушающий доступ к первому элементу
	Items() []ElType       // содержимое очереди от головы к хвосту
}
