package jsonbalance

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

func Balanse_for_json(jsonStr []byte, name string) {
    var data map[string]map[string]map[string]string
    err := json.Unmarshal([]byte(jsonStr), &data)
    if err != nil {
        panic(err)
    }

    file, err := os.Create(name)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    writer.Write([]string{"Code", "Year", "Key", "Value"})

    for code, years := range data {
        for year, keys := range years {
            for key, value := range keys {
                row := []string{code, year, key, value}
                writer.Write(row)
            }
        }
    }

    fmt.Println("JSON успешно преобразован и сохранён в output.csv")
}
