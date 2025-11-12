package jsonconv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type IPDetails struct {
	Сумма      float64  `json:"Сумма"`
	Количество int      `json:"Количество"`
	ИП         []string `json:"ИП"`
}

type StatusData map[string]IPDetails

type YearData map[string]json.RawMessage

type PersonData map[string]YearData

type APIResponse map[string]PersonData

type IPRecord struct {
	ID            string
	Year          string
	Status        string
	VziskanieType string
	Sum           float64
	Count         int
	IPNumber      string
}

// Parser - основной парсер
type Parser struct {
	records []IPRecord
}

// NewParser создает новый парсер
func NewParser() *Parser {
	return &Parser{
		records: make([]IPRecord, 0),
	}
}

func (p *Parser) ParseJSON(jsonData []byte) error {
	var apiResponse APIResponse

	if err := json.Unmarshal(jsonData, &apiResponse); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	p.records = make([]IPRecord, 0)

	for personID, yearsData := range apiResponse {
		for year, statusesData := range yearsData {
			for status, rawStatusData := range statusesData {
				if len(rawStatusData) == 0 || string(rawStatusData) == "[]" {
					continue
				}

				var typesData StatusData
				if err := json.Unmarshal(rawStatusData, &typesData); err != nil {
					continue
				}

				for vziskanieType, details := range typesData {
					for _, ipNumber := range details.ИП {
						record := IPRecord{
							ID:            personID,
							Year:          year,
							Status:        status,
							VziskanieType: vziskanieType,
							Sum:           details.Сумма,
							Count:         details.Количество,
							IPNumber:      ipNumber,
						}
						p.records = append(p.records, record)
					}
				}
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

	headers := []string{"ID", "Год", "Статус", "Тип взыскания", "Сумма", "Количество", "Номер ИП"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("ошибка записи заголовков: %w", err)
	}

	for _, record := range p.records {
		row := []string{
			record.ID,
			record.Year,
			record.Status,
			record.VziskanieType,
			fmt.Sprintf("%.2f", record.Sum),
			strconv.Itoa(record.Count),
			record.IPNumber,
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

	result += "ID,Год,Статус,Тип взыскания,Сумма,Количество,Номер ИП\n"

	for _, record := range p.records {
		result += fmt.Sprintf("%s,%s,%s,%s,%.2f,%d,%s\n",
			record.ID,
			record.Year,
			record.Status,
			record.VziskanieType,
			record.Sum,
			record.Count,
			record.IPNumber,
		)
	}

	return result, nil
}

// Stats - статистика по данным
type Stats struct {
	TotalRecords   int
	TotalSum       float64
	Years          []string
	Statuses       []string
	VziskanieTypes int
	UniquePersons  int
}

func (p *Parser) GetStats() Stats {
	if len(p.records) == 0 {
		return Stats{}
	}

	var totalSum float64
	yearsMap := make(map[string]bool)
	statusesMap := make(map[string]bool)
	typesMap := make(map[string]bool)
	personsMap := make(map[string]bool)

	for _, record := range p.records {
		totalSum += record.Sum
		yearsMap[record.Year] = true
		statusesMap[record.Status] = true
		typesMap[record.VziskanieType] = true
		personsMap[record.ID] = true
	}

	years := make([]string, 0, len(yearsMap))
	for year := range yearsMap {
		years = append(years, year)
	}
	sort.Strings(years)

	statuses := make([]string, 0, len(statusesMap))
	for status := range statusesMap {
		statuses = append(statuses, status)
	}
	sort.Strings(statuses)

	return Stats{
		TotalRecords:   len(p.records),
		TotalSum:       totalSum,
		Years:          years,
		Statuses:       statuses,
		VziskanieTypes: len(typesMap),
		UniquePersons:  len(personsMap),
	}
}

func (p *Parser) GetRecords() []IPRecord {
	return p.records
}

func (p *Parser) FilterByYear(year string) []IPRecord {
	filtered := make([]IPRecord, 0)
	for _, record := range p.records {
		if record.Year == year {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

// FilterByStatus фильтрует записи по статусу
func (p *Parser) FilterByStatus(status string) []IPRecord {
	filtered := make([]IPRecord, 0)
	for _, record := range p.records {
		if record.Status == status {
			filtered = append(filtered, record)
		}
	}
	return filtered
}
