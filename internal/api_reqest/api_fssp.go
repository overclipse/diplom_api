package api_reqest

import (
	"bytes"
	"diplom_api/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tuan78/jsonconv"
	"github.com/yukithm/json2csv"
)

func Fsssp() {
	finalURL := "https://api.damia.ru/fssp/isps?inn=7728551528&format=1&key=2268a80e1f11a48f8657c69c71f1c00d41be9219"
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
	if err != nil {
		log.Fatal(err)
	}
	w := json2csv.NewCSVWriter(os.Stdout)
}