package errdiplom

import "errors"

var (
	Errfilecsv  = errors.New("Ошибка записи csv")
	Errcorparse = errors.New("Ошибка парсинга byte")
)
