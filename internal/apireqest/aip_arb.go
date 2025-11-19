package apireqest

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"diplom_api/internal/jsonconv/arbitrage"
)

func Arb(url string) {
	if url == "" {
		log.Println("Arb: пустой url")
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println("http.Get error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("read body error:", err)
		return
	}
	fmt.Println("Значение arb получено")

	parser := arbitrage.NewParser()

	if err := parser.ParseJSON(body); err != nil {
		log.Println("ParseJSON error:", err)
		return
	}

	if err := parser.ToCSV("arbitrage.csv"); err != nil {
		log.Println("ToCSV error:", err)
		return
	}
}

