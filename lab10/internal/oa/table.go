package oa

type slot struct {
	used bool
	key  string
}

type Table struct {
	slots []slot
}

func NewTable(m int) *Table {
	return &Table{slots: make([]slot, m)}
}

func (t *Table) Size() int { return len(t.slots) }

// Insert измеряет длину пути поиска (число проверенных ячеек) при вставке.
// Возвращает (probes, ok). ok=false если не нашли место за M проб.
func (t *Table) Insert(key string, h0 int, kind ProbeKind, pp ProbeParams) (int, bool) {
	m := t.Size()

	for i := 0; i < m; i++ {
		addr := ProbeAddress(h0, i, kind, pp)
		s := &t.slots[addr]

		// 1 проба = проверка одной ячейки
		probes := i + 1

		if !s.used {
			s.used = true
			s.key = key
			return probes, true
		}

		// Если тот же ключ уже есть — считаем, что "нашли" (вставка не нужна)
		if s.key == key {
			return probes, true
		}
	}

	return m, false
}
