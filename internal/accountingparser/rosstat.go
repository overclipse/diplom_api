package accountingparser

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type YearData map[string]string

type CompanyData map[string]YearData

type APIResponse map[string]CompanyData

type AccountingRecord struct {
	CompanyINN string  // ИНН компании
	Year       string  // Год отчетности
	LineCode   string  // Код строки (1100, 2100 и т.д.)
	Value      float64 // Значение
}

type Parser struct {
	records []AccountingRecord
}

func NewParser() *Parser {
	return &Parser{
		records: make([]AccountingRecord, 0),
	}
}

func (p *Parser) ParseJSON(jsonData []byte) error {
	var apiResponse APIResponse

	if err := json.Unmarshal(jsonData, &apiResponse); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	p.records = make([]AccountingRecord, 0)

	for companyINN, yearsData := range apiResponse {
		for year, linesData := range yearsData {
			for lineCode, valueStr := range linesData {
				value, err := strconv.ParseFloat(valueStr, 64)
				if err != nil {
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

func (p *Parser) ToCSV(filename string) error {
	if len(p.records) == 0 {
		return fmt.Errorf("нет данных для экспорта")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	file.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ИНН", "Год", "Код строки", "Значение"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("ошибка записи заголовков: %w", err)
	}

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

func (p *Parser) ToCSVString() (string, error) {
	if len(p.records) == 0 {
		return "", fmt.Errorf("нет данных для экспорта")
	}

	var result string

	result += "ИНН,Год,Код строки,Значение\n"

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

type Stats struct {
	TotalRecords    int
	Years           []string
	LineCodes       int
	UniqueCompanies int
}

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

func (p *Parser) GetRecords() []AccountingRecord {
	return p.records
}

func (p *Parser) FilterByYear(year string) []AccountingRecord {
	filtered := make([]AccountingRecord, 0)
	for _, record := range p.records {
		if record.Year == year {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func (p *Parser) FilterByLineCode(lineCode string) []AccountingRecord {
	filtered := make([]AccountingRecord, 0)
	for _, record := range p.records {
		if record.LineCode == lineCode {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func (p *Parser) GetYearData(companyINN, year string) map[string]float64 {
	result := make(map[string]float64)
	for _, record := range p.records {
		if record.CompanyINN == companyINN && record.Year == year {
			result[record.LineCode] = record.Value
		}
	}
	return result
}

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
}

func GetLineCodeDescription(lineCode string) string {
	if desc, ok := LineCodeDescriptions[lineCode]; ok {
		return desc
	}
	return "Неизвестный код"
}
