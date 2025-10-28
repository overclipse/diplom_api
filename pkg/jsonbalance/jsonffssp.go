package jsonbalance

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "os"
    "strconv"
)

// Универсальное преобразование значения в строку
func toString(value interface{}) string {
    switch v := value.(type) {
    case string:
        return v
    case float64:
        // Форматируем float без экспоненты
        return strconv.FormatFloat(v, 'f', -1, 64)
    case int:
        return strconv.Itoa(v)
    default:
        return fmt.Sprintf("%v", v)
    }
}

// Функция конвертации JSON из байтов в CSV файл
func Convert(jsonBytes []byte, outputPath string) error {
    var data map[string]interface{}
    err := json.Unmarshal(jsonBytes, &data)
    if err != nil {
        return err
    }

    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Заголовки
    writer.Write([]string{"Code", "Year", "Status", "Value"})

    for code, yearsRaw := range data {
        years, ok := yearsRaw.(map[string]interface{})
        if !ok {
            continue
        }
        for year, keysRaw := range years {
            keys, ok := keysRaw.(map[string]interface{})
            if !ok {
                continue
            }
            for key, value := range keys {
                strValue := toString(value)
                writer.Write([]string{code, year, key, strValue})
            }
        }
    }

    return nil
}
