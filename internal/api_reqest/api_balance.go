package api_reqest

import (
	"diplom_api/internal/config"
	"diplom_api/pkg/jsonbalance"
	"io"
	"net/http"
	"time"
)

func Api_balance() {

	/* def := config.APIConfig{
	BaseURL: "https://api.damia.ru",
	Method:  "scoring/fincoefs",
	Model:   "_bankrots2016",
	INN:     "6663003127",
	}

	cfg := config.ReadFromConsole(def)
	finalURL := cfg.BuildURL()
	fmt.Println("Final URL:", finalURL) */

	// 1) HTTP-запрос
	finalURL := "https://api.damia.ru/rs/balance?inn=7728551528&key=67266ba78d7779083310826cc491faf420858d1f"
	req, err := http.NewRequest(http.MethodGet, finalURL, nil)
	config.Must(err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "fincoefs-minimal-client")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	config.Must(err)
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		panic("non-2xx status: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	config.Must(err)
	jsonbalance.Balanse_for_json(body, "balance.csv")
}