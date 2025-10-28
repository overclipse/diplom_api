package main

import (
	"diplom_api/pkg/jsonbalance"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Укажите ваш URL API
	apiURL := "https://api.damia.ru/fssp/isps?inn=7728551528&format=1&key=2268a80e1f11a48f8657c69c71f1c00d41be9219"

	fmt.Println("Получаю данные из API...")

	// Получаем данные
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("❌ Ошибка запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Ошибка чтения: %v\n", err)
		return
	}
	csv, err :=  jsonbalance.Unmarshal
	if err != nil {
		fmt.Printf("@w", err)
	}
	fmt.Println(csv)
}
