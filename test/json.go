package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// ТВОЯ структура + json-теги, чтобы надёжно маппилось по ключам из JSON
type Main struct {
	Code     string `json:"Code"`
	Document struct {
		НомерДокумента    string  `json:"НомерДокумента"`
		Дата              string  `json:"Дата"`
		РегНомерСД        string  `json:"РегНомерСД"`
		ВидИсп            string  `json:"ВидИсп"`
		ДатаИсп           string  `json:"ДатаИсп"`
		НомерИсп          string  `json:"НомерИсп"`
		ИННОргИсп         string  `json:"ИННОргИсп"`
		Предмет           string  `json:"Предмет"`
		Должник_НаимФССП  string  `json:"Должник_НаимФССП"`
		Должник_АдресФССП string  `json:"Должник_АдресФССП"`
		Сумма             float64 `json:"Сумма"`
		Остаток           float64 `json:"Остаток"`
		ДепНаим           string  `json:"ДепНаим"`
		ДепАдрес          string  `json:"ДепАдрес"`
		Статус            string  `json:"Статус"`
		ДатаЗаверш        string  `json:"ДатаЗаверш"`
		ПричЗаверш        string  `json:"ПричЗаверш"`
	} `json:"Document"`
}

// Заголовок CSV в фиксированном порядке (как тебе удобно)
var header = []string{
	"Code",
	"НомерДокумента", "Дата", "РегНомерСД", "ВидИсп", "ДатаИсп", "НомерИсп", "ИННОргИсп", "Предмет",
	"Должник_НаимФССП", "Должник_АдресФССП",
	"Сумма", "Остаток",
	"ДепНаим", "ДепАдрес",
	"Статус", "ДатаЗаверш", "ПричЗаверш",
}

// Преобразование []byte JSON -> []byte CSV
func CSVFromMainJSON(data []byte) ([]byte, error) {
	recs, err := parseMain(data)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := writeCSV(&buf, recs); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Пишем сразу в io.Writer (например, в файл)
func writeCSV(w io.Writer, recs []Main) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := cw.Write(header); err != nil {
		return err
	}
	for _, m := range recs {
		row := []string{
			m.Code,
			m.Document.НомерДокумента,
			m.Document.Дата,
			m.Document.РегНомерСД,
			m.Document.ВидИсп,
			m.Document.ДатаИсп,
			m.Document.НомерИсп,
			m.Document.ИННОргИсп,
			m.Document.Предмет,
			m.Document.Должник_НаимФССП,
			m.Document.Должник_АдресФССП,
			floatToStr(m.Document.Сумма),
			floatToStr(m.Document.Остаток),
			m.Document.ДепНаим,
			m.Document.ДепАдрес,
			m.Document.Статус,
			m.Document.ДатаЗаверш,
			m.Document.ПричЗаверш,
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	return cw.Error()
}

func floatToStr(f float64) string {
	// Без лишних нулей: 123.45 / 0 / 1000
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// Пытаемся распарсить либо массив Main, либо одиночный Main
func parseMain(data []byte) ([]Main, error) {
	var arr []Main
	if err := json.Unmarshal(data, &arr); err == nil {
		return arr, nil
	}
	var one Main
	if err := json.Unmarshal(data, &one); err == nil {
		return []Main{one}, nil
	}
	return nil, fmt.Errorf("input is neither Main nor []Main JSON")
}

// Пример использования
