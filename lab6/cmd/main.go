package main

import (
	"bufio"
	"container/heap"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	RecordSize = 200
	KeySize    = 4
	DataSize   = RecordSize - KeySize
	DefaultN   = 10000
)

type Record struct {
	Key  int32
	Data [DataSize]byte
}

func (r *Record) FromBytes(b []byte) {
	r.Key = int32(binary.LittleEndian.Uint32(b[:4]))
	copy(r.Data[:], b[4:])
}

func (r *Record) ToBytes(b []byte) {
	binary.LittleEndian.PutUint32(b[:4], uint32(r.Key))
	copy(b[4:], r.Data[:])
}

func readRecord(rd *bufio.Reader, rec *Record) (bool, error) {
	buf := make([]byte, RecordSize)
	_, err := io.ReadFull(rd, buf)
	if err == io.EOF {
		return false, nil
	}
	if err == io.ErrUnexpectedEOF {
		return false, fmt.Errorf("файл поврежден: размер не кратен %d байтам", RecordSize)
	}
	if err != nil {
		return false, err
	}
	rec.FromBytes(buf)
	return true, nil
}

func writeRecord(w *bufio.Writer, rec *Record) error {
	buf := make([]byte, RecordSize)
	rec.ToBytes(buf)
	_, err := w.Write(buf)
	return err
}

// ---------- Генерация тестового файла ----------
func generateFile(path string, n int, keyMax int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriterSize(f, 1<<20)
	defer w.Flush()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rec := Record{}
	for i := 0; i < n; i++ {
		rec.Key = int32(rng.Intn(keyMax)) // ключи могут повторяться
		for j := 0; j < DataSize; j++ {
			rec.Data[j] = byte(rng.Intn(256))
		}
		if err := writeRecord(w, &rec); err != nil {
			return err
		}
	}
	return nil
}

// ---------- Создание начальных серий ----------
func makeInitialRuns(inputPath, tempDir string, memRecords int) ([]string, error) {
	in, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	rd := bufio.NewReaderSize(in, 1<<20)

	var runs []string
	chunk := make([]Record, 0, memRecords)

	for {
		chunk = chunk[:0]

		for len(chunk) < memRecords {
			var rec Record
			ok, err := readRecord(rd, &rec)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
			chunk = append(chunk, rec)
		}

		if len(chunk) == 0 {
			break
		}

		sort.Slice(chunk, func(i, j int) bool {
			return chunk[i].Key < chunk[j].Key
		})

		runPath := filepath.Join(tempDir, fmt.Sprintf("run_%06d.dat", len(runs)))
		if err := writeRun(runPath, chunk); err != nil {
			return nil, err
		}
		runs = append(runs, runPath)

		if len(chunk) < memRecords {
			break
		}
	}

	return runs, nil
}

func writeRun(path string, records []Record) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriterSize(f, 1<<20)
	defer w.Flush()

	for i := range records {
		if err := writeRecord(w, &records[i]); err != nil {
			return err
		}
	}
	return nil
}

// ---------- k-way merge ----------
type runReader struct {
	f  *os.File
	rd *bufio.Reader
}

type heapItem struct {
	rec   Record
	srcID int
}

type recHeap []heapItem

func (h recHeap) Len() int { return len(h) }
func (h recHeap) Less(i, j int) bool {
	if h[i].rec.Key != h[j].rec.Key {
		return h[i].rec.Key < h[j].rec.Key
	}
	return h[i].srcID < h[j].srcID
}
func (h recHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *recHeap) Push(x any)   { *h = append(*h, x.(heapItem)) }
func (h *recHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func mergeKRuns(runPaths []string, outPath string) error {
	readers := make([]runReader, 0, len(runPaths))
	for _, p := range runPaths {
		f, err := os.Open(p)
		if err != nil {
			for i := range readers {
				readers[i].f.Close()
			}
			return err
		}
		readers = append(readers, runReader{
			f:  f,
			rd: bufio.NewReaderSize(f, 1<<20),
		})
	}
	defer func() {
		for i := range readers {
			readers[i].f.Close()
		}
	}()

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriterSize(out, 1<<20)
	defer w.Flush()

	h := &recHeap{}
	heap.Init(h)

	for i := range readers {
		var rec Record
		ok, err := readRecord(readers[i].rd, &rec)
		if err != nil {
			return err
		}
		if ok {
			heap.Push(h, heapItem{rec: rec, srcID: i})
		}
	}

	for h.Len() > 0 {
		item := heap.Pop(h).(heapItem)
		if err := writeRecord(w, &item.rec); err != nil {
			return err
		}

		var next Record
		ok, err := readRecord(readers[item.srcID].rd, &next)
		if err != nil {
			return err
		}
		if ok {
			heap.Push(h, heapItem{rec: next, srcID: item.srcID})
		}
	}

	return nil
}

// ---------- Сортировка ----------
type SortStats struct {
	InitialRuns int
	MergePasses int
	TotalPasses int
	Duration    time.Duration
}

func externalBalancedMultiwayMergeSort(inPath, outPath string, k, memRecords int) (SortStats, error) {
	if k < 2 {
		return SortStats{}, fmt.Errorf("k должно быть >= 2")
	}
	if memRecords < 1 {
		return SortStats{}, fmt.Errorf("mem должно быть >= 1")
	}

	start := time.Now()

	tempDir, err := os.MkdirTemp("", "kway_sort_*")
	if err != nil {
		return SortStats{}, err
	}
	defer os.RemoveAll(tempDir)

	runs, err := makeInitialRuns(inPath, tempDir, memRecords)
	if err != nil {
		return SortStats{}, err
	}

	stats := SortStats{InitialRuns: len(runs)}

	if len(runs) == 0 {
		_, err := os.Create(outPath)
		stats.Duration = time.Since(start)
		stats.TotalPasses = 1
		return stats, err
	}

	if len(runs) == 1 {
		if err := copyFile(runs[0], outPath); err != nil {
			return stats, err
		}
		stats.MergePasses = 0
		stats.TotalPasses = 1
		stats.Duration = time.Since(start)
		return stats, nil
	}

	pass := 0
	for len(runs) > 1 {
		pass++
		var newRuns []string

		for i := 0; i < len(runs); i += k {
			end := i + k
			if end > len(runs) {
				end = len(runs)
			}
			group := runs[i:end]

			mergedPath := filepath.Join(tempDir, fmt.Sprintf("merge_p%03d_%06d.dat", pass, len(newRuns)))
			if err := mergeKRuns(group, mergedPath); err != nil {
				return stats, err
			}
			newRuns = append(newRuns, mergedPath)

			for _, p := range group {
				_ = os.Remove(p)
			}
		}

		runs = newRuns
	}

	if err := copyFile(runs[0], outPath); err != nil {
		return stats, err
	}

	stats.MergePasses = pass
	stats.TotalPasses = pass + 1
	stats.Duration = time.Since(start)
	return stats, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// ---------- Проверка ----------
func checkSorted(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	rd := bufio.NewReaderSize(f, 1<<20)
	var prev *int32
	for {
		var rec Record
		ok, err := readRecord(rd, &rec)
		if err != nil {
			return err
		}
		if !ok {
			break
		}
		if prev != nil && rec.Key < *prev {
			return fmt.Errorf("файл НЕ отсортирован: встретили %d после %d", rec.Key, *prev)
		}
		v := rec.Key
		prev = &v
	}
	return nil
}

// ---------- Удобный консольный UI ----------
func runUI() {
	in := bufio.NewReader(os.Stdin)

	fmt.Println("==============================================")
	fmt.Println(" Внешняя сортировка: сбалансированное k-way слияние")
	fmt.Println(" Формат записи: 200 байт (первые 4 байта — int32 ключ)")
	fmt.Println("==============================================")

	// значения по умолчанию
	defIn := "input.dat"
	defOut := "sorted.dat"
	defN := DefaultN
	defKeyMax := 10000
	defK := 4
	defMem := 1000

	for {
		fmt.Println("\nМеню:")
		fmt.Println("  1) Сгенерировать входной файл")
		fmt.Println("  2) Отсортировать файл (k-way merge)")
		fmt.Println("  3) Замерить время/проходы для нескольких k (benchmark)")
		fmt.Println("  4) Проверить, что файл отсортирован")
		fmt.Println("  0) Выход")

		choice := promptInt(in, "Выбери пункт", 0)
		switch choice {
		case 1:
			path := promptString(in, "Путь к файлу", defIn)
			n := promptInt(in, "Количество записей", defN)
			keyMax := promptInt(in, "Максимальное значение ключа (0..keyMax-1)", defKeyMax)

			if err := generateFile(path, n, keyMax); err != nil {
				fmt.Println("Ошибка генерации:", err)
			} else {
				fmt.Println("Готово! Сгенерирован файл:", path)
			}

		case 2:
			inPath := promptString(in, "Входной файл", defIn)
			outPath := promptString(in, "Выходной файл", defOut)
			k := promptInt(in, "k (кол-во путей слияния, >=2)", defK)
			mem := promptInt(in, "Сколько записей сортировать в памяти (размер серии)", defMem)
			doCheck := promptBool(in, "Проверить результат (y/n)", true)

			stats, err := externalBalancedMultiwayMergeSort(inPath, outPath, k, mem)
			if err != nil {
				fmt.Println("Ошибка сортировки:", err)
				break
			}

			printStats(stats, k, mem, outPath)

			if doCheck {
				if err := checkSorted(outPath); err != nil {
					fmt.Println("Проверка:", err)
				} else {
					fmt.Println("Проверка: файл отсортирован ✅")
				}
			}

		case 3:
			inPath := promptString(in, "Входной файл (один и тот же для всех k)", defIn)
			outBase := promptString(in, "База имени выходного файла (например sorted)", "sorted")
			mem := promptInt(in, "Сколько записей сортировать в памяти", defMem)
			kList := promptIntList(in, "Список k через пробел", []int{2, 4, 8, 16})
			doCheck := promptBool(in, "Проверять каждый результат (y/n)", false)

			fmt.Println("\nРезультаты:")
			fmt.Println("k\tнач.серии\tпрох.слияния\tвсего\tвремя")

			for _, k := range kList {
				outPath := fmt.Sprintf("%s_k%d.dat", outBase, k)
				stats, err := externalBalancedMultiwayMergeSort(inPath, outPath, k, mem)
				if err != nil {
					fmt.Printf("%d\t-\t-\t-\tОШИБКА: %v\n", k, err)
					continue
				}

				fmt.Printf("%d\t%d\t\t%d\t\t%d\t%s\n",
					k, stats.InitialRuns, stats.MergePasses, stats.TotalPasses, stats.Duration)

				if doCheck {
					if err := checkSorted(outPath); err != nil {
						fmt.Println("  проверка:", err)
					} else {
						fmt.Println("  проверка: OK ✅")
					}
				}
			}

		case 4:
			path := promptString(in, "Файл для проверки", defOut)
			if err := checkSorted(path); err != nil {
				fmt.Println("Проверка:", err)
			} else {
				fmt.Println("Проверка: файл отсортирован ✅")
			}

		case 0:
			fmt.Println("Пока!")
			return

		default:
			fmt.Println("Неизвестный пункт меню.")
		}
	}
}

func printStats(stats SortStats, k, mem int, outPath string) {
	fmt.Println("\nГотово!")
	fmt.Println("Выходной файл:", outPath)
	fmt.Printf("Параметры: k=%d, mem=%d записей\n", k, mem)
	fmt.Printf("Начальных серий: %d\n", stats.InitialRuns)
	fmt.Printf("Проходов слияния: %d\n", stats.MergePasses)
	fmt.Printf("Всего проходов (включая создание серий): %d\n", stats.TotalPasses)
	fmt.Printf("Время сортировки: %s\n", stats.Duration)
}

func promptString(r *bufio.Reader, msg, def string) string {
	for {
		if def != "" {
			fmt.Printf("%s [%s]: ", msg, def)
		} else {
			fmt.Printf("%s: ", msg)
		}
		line, _ := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			return def
		}
		return line
	}
}

func promptInt(r *bufio.Reader, msg string, def int) int {
	for {
		fmt.Printf("%s [%d]: ", msg, def)
		line, _ := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			return def
		}
		v, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("Введите целое число.")
			continue
		}
		return v
	}
}

func promptBool(r *bufio.Reader, msg string, def bool) bool {
	defStr := "y"
	if !def {
		defStr = "n"
	}
	for {
		fmt.Printf("%s [%s]: ", msg, defStr)
		line, _ := r.ReadString('\n')
		line = strings.TrimSpace(strings.ToLower(line))
		if line == "" {
			return def
		}
		if line == "y" || line == "yes" || line == "д" || line == "да" {
			return true
		}
		if line == "n" || line == "no" || line == "н" || line == "нет" {
			return false
		}
		fmt.Println("Введите y/n (да/нет).")
	}
}

func promptIntList(r *bufio.Reader, msg string, def []int) []int {
	defStrs := make([]string, 0, len(def))
	for _, v := range def {
		defStrs = append(defStrs, strconv.Itoa(v))
	}
	defStr := strings.Join(defStrs, " ")

	for {
		fmt.Printf("%s [%s]: ", msg, defStr)
		line, _ := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			return def
		}
		parts := strings.Fields(line)
		var res []int
		ok := true
		for _, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil || v < 2 {
				ok = false
				break
			}
			res = append(res, v)
		}
		if !ok || len(res) == 0 {
			fmt.Println("Введите числа >=2 через пробел, например: 2 4 8 16")
			continue
		}
		return res
	}
}

// ---------- main ----------
func main() {
	// Если запуск без аргументов — сразу UI (удобно для лабораторной)
	if len(os.Args) == 1 {
		runUI()
		return
	}

	// Флаги (как было раньше) + возможность принудительно UI
	inPath := flag.String("in", "input.dat", "входной файл")
	outPath := flag.String("out", "sorted.dat", "выходной файл")
	k := flag.Int("k", 4, "количество путей слияния (k)")
	mem := flag.Int("mem", 1000, "сколько записей сортируем в памяти для одной серии")
	gen := flag.Bool("gen", false, "сгенерировать входной файл (по умолчанию 10000 записей)")
	n := flag.Int("n", DefaultN, "количество записей при генерации")
	keyMax := flag.Int("keymax", 10000, "максимальное значение ключа при генерации")
	check := flag.Bool("check", true, "проверить, что выходной файл отсортирован")
	ui := flag.Bool("ui", false, "открыть интерактивное меню")
	flag.Parse()

	if *ui {
		runUI()
		return
	}

	if *gen {
		if err := generateFile(*inPath, *n, *keyMax); err != nil {
			fmt.Println("Ошибка генерации:", err)
			os.Exit(1)
		}
		fmt.Println("Сгенерирован файл:", *inPath)
	}

	stats, err := externalBalancedMultiwayMergeSort(*inPath, *outPath, *k, *mem)
	if err != nil {
		fmt.Println("Ошибка сортировки:", err)
		os.Exit(1)
	}

	printStats(stats, *k, *mem, *outPath)

	if *check {
		if err := checkSorted(*outPath); err != nil {
			fmt.Println("Проверка:", err)
			os.Exit(1)
		}
		fmt.Println("Проверка: файл отсортирован ✅")
	}
}
