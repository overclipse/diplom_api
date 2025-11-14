package apireqest

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"diplom_api/errdiplom"
	"diplom_api/internal/jsonconv/fssp"
)

func Fsssp(url string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "fincoefs-minimal-client")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		panic("non-2xx status: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("fssp значение получено")
	parser := fssp.NewParser()
	err = parser.ParseJSON(body)
	if err != nil {
		fmt.Errorf("%w: %v", errdiplom.Errcorparse.Error(), err)
	}
	err = parser.ToCSV("out.csv")
	if err != nil {
		fmt.Errorf(errdiplom.Errfilecsv.Error())
	}
}
