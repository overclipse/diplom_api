package config

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
)

// APIConfig — одна структура, из которой собираем финальный URL
type APIConfig struct {
	BaseURL string `json:"base_url"`  // напр.: https://api.damia.ru
	Method  string `json:"method"`    // напр.: scoring/fincoefs или scoring/score
	Model   string `json:"model"`
	INN     string `json:"inn"`
	Key     string `json:"key"`
}

// BuildURL — склейка base + method + query
func (c APIConfig) BuildURL() string {
	if c.BaseURL == "" || c.Method == "" {
		panic("BaseURL and Method are required")
	}
	u, err := url.Parse(c.BaseURL)
	Must(err)

	// аккуратно со слэшами
	u.Path = strings.TrimSuffix(u.Path, "/")
	method := strings.TrimPrefix(c.Method, "/")
	u.Path = path.Join(u.Path, method)

	q := u.Query()
	if c.Model != "" { q.Set("model", c.Model) }
	if c.INN   != "" { q.Set("inn", c.INN) }
	if c.Key   != "" { q.Set("key", c.Key) }
	u.RawQuery = q.Encode()
	return u.String()
}

// ReadFromConsole запрашивает поля у пользователя и возвращает собранный APIConfig.
// В качестве дефолтов использует значения из аргумента def.
func ReadFromConsole(def APIConfig) APIConfig {
	r := bufio.NewReader(os.Stdin)

	read := func(prompt, d string) string {
		if d != "" {
			fmt.Printf("%s [%s]: ", prompt, d)
		} else {
			fmt.Printf("%s: ", prompt)
		}
		s, err := r.ReadString('\n')
		Must(err)
		s = strings.TrimSpace(s)
		if s == "" {
			return d
		}
		return s
	}

	cfg := APIConfig{
		BaseURL: read("Base URL", def.BaseURL),          // напр.: https://api.damia.ru
		Method:  read("Method (endpoint)", def.Method),  // напр.: scoring/fincoefs или scoring/score
		Model:   read("model", def.Model),
		INN:     read("inn", def.INN),
		Key:     read("key", def.Key),
	}

	return cfg
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
