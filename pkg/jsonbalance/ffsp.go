package jsonbalance

import (
"encoding/json"
"fmt"
)

// ----- Листовая структура -----

type Details struct {
	// Поля в твоём JSON называются по-русски — так и маппим:
	Amount     float64  `json:"Сумма"`      // может быть 0 или дробное
	Count      int      `json:"Количество"` // целое
	Executions []string `json:"ИП"`         // список номеров ИП
}

// Categories — это карта "Категория" -> Details,
// но в источнике раздел может прийти как {} / [] / null.
// Делаем кастомный Unmarshal, чтобы нормализовать к пустой map при []/null.
type Categories map[string]Details

func (c *Categories) UnmarshalJSON(b []byte) error {
	// null или пустой массив трактуем как пустой раздел
	if string(b) == "null" || string(b) == "[]" {
		*c = make(Categories)
		return nil
	}

	// Пытаемся как объект
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err == nil {
		out := make(Categories, len(raw))
		for k, v := range raw {
			var d Details
			if err := json.Unmarshal(v, &d); err != nil {
				return fmt.Errorf("category %q decode: %w", k, err)
			}
			out[k] = d
		}
		*c = out
		return nil
	}

	// Если это массив (например []), но не null — считаем пустым
	var arr []any
	if err := json.Unmarshal(b, &arr); err == nil {
		*c = make(Categories)
		return nil
	}

	// Иное — ошибка формата
	return fmt.Errorf("categories: unsupported JSON: %s", string(b))
}

// ----- Блок года -----

type YearBlock struct {
	Completed  Categories `json:"Завершено"`   // map[Категория]Details или пусто
	Repaid     Categories `json:"Погашено"`    // map[Категория]Details или пусто
	Unfinished Categories `json:"Не завершено"`// map[Категория]Details или пусто
}

// Dataset — INN -> Year -> YearBlock
type Dataset map[string]map[string]YearBlock

// ParseFSSPJSON парсит []byte от HTTP-ответа в строго типизированную структуру.
func ParseFSSPJSON(b []byte) (Dataset, error) {
	var ds Dataset
	if err := json.Unmarshal(b, &ds); err != nil {
		return nil, fmt.Errorf("unmarshal dataset: %w", err)
	}
	// На этом этапе все разделы уже нормализованы ({} вместо []/null).
	return ds, nil
}
