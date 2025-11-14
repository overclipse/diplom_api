package apireqest

import (
	"diplom_api/internal/jsonconv/rosstat"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Rosstat(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("значение росстат полученно")

	parser := rosstat.NewParser()

	if err := parser.ParseJSON(body); err != nil {
		log.Fatal(err)
	}

	if err := parser.ToCSV("accounting.csv"); err != nil {
		log.Fatal(err)
	}
	stats := parser.GetStats()
	fmt.Printf("Всего записей: %d\n", stats.TotalRecords)
	fmt.Printf("Годы: %v\n", stats.Years)
	fmt.Printf("Уникальных кодов: %d\n", stats.LineCodes)
	fmt.Printf("Компаний: %d\n", stats.UniqueCompanies)

	// Получить данные за конкретный год
	yearData := parser.GetYearData("7713076301", "2024")
	fmt.Printf("Выручка 2024: %.0f\n", yearData["2100"])
	fmt.Printf("Чистая прибыль 2024: %.0f\n", yearData["2500"])

	// Фильтрация по коду строки (например, выручка)
	revenues := parser.FilterByLineCode("2100")
	fmt.Printf("Записей о выручке: %d\n", len(revenues))
}
