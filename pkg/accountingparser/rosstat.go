package accountingparser

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// YearData содержит данные по кодам строк для года
type YearData map[string]string

// CompanyData содержит данные по годам для компании
type CompanyData map[string]YearData

// APIResponse - корневая структура ответа API
type APIResponse map[string]CompanyData

// AccountingRecord - запись для CSV
type AccountingRecord struct {
	CompanyINN string  // ИНН компании
	Year       string  // Год отчетности
	LineCode   string  // Код строки (1100, 2100 и т.д.)
	Value      float64 // Значение
}

// Parser - основной парсер бухгалтерской отчетности
type Parser struct {
	records []AccountingRecord
}

// NewParser создает новый парсер
func NewParser() *Parser {
	return &Parser{
		records: make([]AccountingRecord, 0),
	}
}

// ParseJSON парсит JSON body в структуры
func (p *Parser) ParseJSON(jsonData []byte) error {
	var apiResponse APIResponse

	if err := json.Unmarshal(jsonData, &apiResponse); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	p.records = make([]AccountingRecord, 0)

	// Проходим по всем компаниям (ИНН)
	for companyINN, yearsData := range apiResponse {
		// Проходим по всем годам
		for year, linesData := range yearsData {
			// Проходим по всем кодам строк
			for lineCode, valueStr := range linesData {
				// Конвертируем строку в число
				value, err := strconv.ParseFloat(valueStr, 64)
				if err != nil {
					// Если не удалось распарсить - пропускаем
					continue
				}

				record := AccountingRecord{
					CompanyINN: companyINN,
					Year:       year,
					LineCode:   lineCode,
					Value:      value,
				}
				p.records = append(p.records, record)
			}
		}
	}

	return nil
}

// ToCSV экспортирует данные в CSV
func (p *Parser) ToCSV(filename string) error {
	if len(p.records) == 0 {
		return fmt.Errorf("нет данных для экспорта")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	// UTF-8 BOM для корректного отображения в Excel
	file.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Заголовки
	headers := []string{"ИНН", "Год", "Код строки", "Значение"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("ошибка записи заголовков: %w", err)
	}

	// Данные
	for _, record := range p.records {
		row := []string{
			record.CompanyINN,
			record.Year,
			record.LineCode,
			fmt.Sprintf("%.0f", record.Value), // Без десятичных знаков
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
	}

	return nil
}

// ToCSVString возвращает CSV как строку
func (p *Parser) ToCSVString() (string, error) {
	if len(p.records) == 0 {
		return "", fmt.Errorf("нет данных для экспорта")
	}

	var result string

	// Заголовки
	result += "ИНН,Год,Код строки,Значение\n"

	// Данные
	for _, record := range p.records {
		result += fmt.Sprintf("%s,%s,%s,%.0f\n",
			record.CompanyINN,
			record.Year,
			record.LineCode,
			record.Value,
		)
	}

	return result, nil
}

// Stats - статистика по данным
type Stats struct {
	TotalRecords    int
	Years           []string
	LineCodes       int
	UniqueCompanies int
}

// GetStats возвращает статистику
func (p *Parser) GetStats() Stats {
	if len(p.records) == 0 {
		return Stats{}
	}

	yearsMap := make(map[string]bool)
	lineCodesMap := make(map[string]bool)
	companiesMap := make(map[string]bool)

	for _, record := range p.records {
		yearsMap[record.Year] = true
		lineCodesMap[record.LineCode] = true
		companiesMap[record.CompanyINN] = true
	}

	// Конвертируем map в slice
	years := make([]string, 0, len(yearsMap))
	for year := range yearsMap {
		years = append(years, year)
	}
	sort.Strings(years)

	return Stats{
		TotalRecords:    len(p.records),
		Years:           years,
		LineCodes:       len(lineCodesMap),
		UniqueCompanies: len(companiesMap),
	}
}

// GetRecords возвращает все записи
func (p *Parser) GetRecords() []AccountingRecord {
	return p.records
}

// FilterByYear фильтрует записи по году
func (p *Parser) FilterByYear(year string) []AccountingRecord {
	filtered := make([]AccountingRecord, 0)
	for _, record := range p.records {
		if record.Year == year {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

// FilterByLineCode фильтрует записи по коду строки
func (p *Parser) FilterByLineCode(lineCode string) []AccountingRecord {
	filtered := make([]AccountingRecord, 0)
	for _, record := range p.records {
		if record.LineCode == lineCode {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

// GetYearData возвращает все данные за конкретный год для компании
func (p *Parser) GetYearData(companyINN, year string) map[string]float64 {
	result := make(map[string]float64)
	for _, record := range p.records {
		if record.CompanyINN == companyINN && record.Year == year {
			result[record.LineCode] = record.Value
		}
	}
	return result
}

// LineCodeDescriptions - справочник описаний кодов строк
var LineCodeDescriptions = map[string]string{
	"1100": "Нематериальные активы",
	"1110": "Результаты исследований и разработок",
	"1150": "Основные средства",
	"1160": "Доходные вложения в материальные ценности",
	"1170": "Финансовые вложения",
	"1180": "Отложенные налоговые активы",
	"1190": "Прочие внеоборотные активы",
	"1200": "Запасы",
	"1210": "НДС по приобретенным ценностям",
	"1220": "Дебиторская задолженность",
	"1230": "Финансовые вложения (за исключением денежных эквивалентов)",
	"1240": "Денежные средства и денежные эквиваленты",
	"1250": "Прочие оборотные активы",
	"1260": "Итого оборотные активы",
	"1300": "Капитал и резервы",
	"1400": "Долгосрочные обязательства",
	"1500": "Краткосрочные обязательства",
	"1600": "БАЛАНС (актив)",
	"1700": "БАЛАНС (пассив)",
	"2100": "Выручка",
	"2110": "Себестоимость продаж",
	"2200": "Валовая прибыль (убыток)",
	"2300": "Прибыль (убыток) от продаж",
	"2400": "Прибыль (убыток) до налогообложения",
	"2500": "Чистая прибыль (убыток)",
	// Добавьте другие коды по необходимости
}

func GetLineCodeDescription(lineCode string) string {
	if desc, ok := LineCodeDescriptions[lineCode]; ok {
		return desc
	}
	return "Неизвестный код"
}
