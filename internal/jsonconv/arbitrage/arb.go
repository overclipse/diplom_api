package arbitrage
import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// DecisionDetails содержит детали решения суда
type DecisionDetails struct {
	Сумма      float64 `json:"Сумма"`
	Количество int     `json:"Количество"`
}

// DecisionsData содержит данные по конкретным решениям
type DecisionsData map[string]DecisionDetails

// DecisionTypeData содержит данные по типам решений
type DecisionTypeData map[string]json.RawMessage

// YearData содержит данные по годам
type YearData map[string]DecisionTypeData

// RoleData содержит данные по ролям (Истец, Ответчик, ТретьеЛицо)
type RoleData map[string]YearData

// APIResponse - корневая структура ответа API
type APIResponse struct {
	Result RoleData `json:"result"`
}

// CourtRecord - запись для CSV
type CourtRecord struct {
	Role         string  // Роль: Истец, Ответчик, ТретьеЛицо
	Year         string  // Год
	DecisionType string  // Тип решения: РешенияПерв, РешенияАпп, РешенияКасс, РешенияНадз, Итого
	Decision     string  // Конкретное решение суда
	Sum          float64 // Сумма
	Count        int     // Количество
}

// Parser - основной парсер судебных дел
type Parser struct {
	records []CourtRecord
}

// NewParser создает новый парсер
func NewParser() *Parser {
	return &Parser{
		records: make([]CourtRecord, 0),
	}
}

// ParseJSON парсит JSON body в структуры
func (p *Parser) ParseJSON(jsonData []byte) error {
	var apiResponse APIResponse

	if err := json.Unmarshal(jsonData, &apiResponse); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	p.records = make([]CourtRecord, 0)

	// Проходим по всем ролям (Истец, Ответчик, ТретьеЛицо)
	for role, yearsData := range apiResponse.Result {
		// Проходим по всем годам
		for year, decisionTypesData := range yearsData {
			// Проходим по всем типам решений
			for decisionType, rawDecisionsData := range decisionTypesData {
				// Проверяем, что это не пустой массив
				if len(rawDecisionsData) == 0 || string(rawDecisionsData) == "[]" || string(rawDecisionsData) == "null" {
					continue
				}

				// Проверяем что это объект
				if len(rawDecisionsData) > 0 && rawDecisionsData[0] != '{' {
					continue
				}

				// Пробуем распарсить как DecisionsData
				var decisionsData DecisionsData
				if err := json.Unmarshal(rawDecisionsData, &decisionsData); err != nil {
					// Если не получилось, пропускаем
					continue
				}

				// Проходим по всем конкретным решениям
				for decision, details := range decisionsData {
					record := CourtRecord{
						Role:         role,
						Year:         year,
						DecisionType: decisionType,
						Decision:     decision,
						Sum:          details.Сумма,
						Count:        details.Количество,
					}
					p.records = append(p.records, record)
				}
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
	headers := []string{"Роль", "Год", "Тип решения", "Решение", "Сумма", "Количество"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("ошибка записи заголовков: %w", err)
	}

	// Данные
	for _, record := range p.records {
		row := []string{
			record.Role,
			record.Year,
			record.DecisionType,
			record.Decision,
			fmt.Sprintf("%.2f", record.Sum),
			fmt.Sprintf("%d", record.Count),
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
	result += "Роль,Год,Тип решения,Решение,Сумма,Количество\n"

	// Данные
	for _, record := range p.records {
		result += fmt.Sprintf("%s,%s,%s,%s,%.2f,%d\n",
			record.Role,
			record.Year,
			record.DecisionType,
			record.Decision,
			record.Sum,
			record.Count,
		)
	}

	return result, nil
}

// Stats - статистика по данным
type Stats struct {
	TotalRecords   int
	TotalSum       float64
	TotalCount     int
	Years          []string
	Roles          []string
	DecisionTypes  int
}

// GetStats возвращает статистику
func (p *Parser) GetStats() Stats {
	if len(p.records) == 0 {
		return Stats{}
	}

	var totalSum float64
	var totalCount int
	yearsMap := make(map[string]bool)
	rolesMap := make(map[string]bool)
	decisionTypesMap := make(map[string]bool)

	for _, record := range p.records {
		totalSum += record.Sum
		totalCount += record.Count
		yearsMap[record.Year] = true
		rolesMap[record.Role] = true
		decisionTypesMap[record.DecisionType] = true
	}

	// Конвертируем map в slice
	years := make([]string, 0, len(yearsMap))
	for year := range yearsMap {
		years = append(years, year)
	}
	sort.Strings(years)

	roles := make([]string, 0, len(rolesMap))
	for role := range rolesMap {
		roles = append(roles, role)
	}
	sort.Strings(roles)

	return Stats{
		TotalRecords:  len(p.records),
		TotalSum:      totalSum,
		TotalCount:    totalCount,
		Years:         years,
		Roles:         roles,
		DecisionTypes: len(decisionTypesMap),
	}
}

// GetRecords возвращает все записи
func (p *Parser) GetRecords() []CourtRecord {
	return p.records
}

// FilterByRole фильтрует записи по роли
func (p *Parser) FilterByRole(role string) []CourtRecord {
	filtered := make([]CourtRecord, 0)
	for _, record := range p.records {
		if record.Role == role {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

// FilterByYear фильтрует записи по году
func (p *Parser) FilterByYear(year string) []CourtRecord {
	filtered := make([]CourtRecord, 0)
	for _, record := range p.records {
		if record.Year == year {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

// FilterByDecisionType фильтрует записи по типу решения
func (p *Parser) FilterByDecisionType(decisionType string) []CourtRecord {
	filtered := make([]CourtRecord, 0)
	for _, record := range p.records {
		if record.DecisionType == decisionType {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

// GetSummary возвращает сводку по роли и году
func (p *Parser) GetSummary(role, year string) map[string]float64 {
	summary := make(map[string]float64)
	
	for _, record := range p.records {
		if record.Role == role && record.Year == year {
			summary[record.DecisionType] += record.Sum
		}
	}
	
	return summary
}

// DecisionTypeDescriptions - справочник типов решений
var DecisionTypeDescriptions = map[string]string{
	"РешенияПерв": "Решения первой инстанции",
	"РешенияАпп":  "Решения апелляционной инстанции",
	"РешенияКасс": "Решения кассационной инстанции",
	"РешенияНадз": "Решения надзорной инстанции",
	"Итого":       "Итоговые данные",
}

// GetDecisionTypeDescription возвращает описание типа решения
func GetDecisionTypeDescription(decisionType string) string {
	if desc, ok := DecisionTypeDescriptions[decisionType]; ok {
		return desc
	}
	return "Неизвестный тип"
}