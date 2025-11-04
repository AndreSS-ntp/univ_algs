package trycatch

import "fmt"

// Do выполняет action и перехватывает любую панику, возвращая её в виде ошибки.
// Такой подход позволяет имитировать конструкцию try/catch, которая отсутствует
// в Go: потенциально аварийный код вызывается внутри Do, а вызывающий код
// получает ошибку вместо падения всей программы.
func Do(action func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
				return
			}
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	action()
	return nil
}
