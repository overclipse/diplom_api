package ui

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/charmbracelet/huh"
)

const (
	defaultFsspEndpoint    = "https://api.damia.ru/fssp/isps"
	defaultRosstatEndpoint = "https://api.damia.ru/rs/balance"
	defaultArbEndpoint     = "https://api.damia.ru/arb/dela"
)

// ErrMenuAborted is returned when the operator aborts the menu.
var ErrMenuAborted = errors.New("request configuration canceled by user")

// RequestConfig keeps the parameters the APIs require.
type RequestConfig struct {
	Fssp    FsspConfig
	Rosstat RosstatConfig
	Arb     ArbConfig
}

// FsspConfig describes the FSSP API request.
type FsspConfig struct {
	Endpoint string
	Inn      string
	Format   string
	Key      string
}

// RosstatConfig describes the Rosstat API request.
type RosstatConfig struct {
	Endpoint string
	Inn      string
	Key      string
}

// ArbConfig describes the Arbitration cases API request.
type ArbConfig struct {
	Endpoint string
	Query    string
	Format   string
	Key      string
}

// RequestURLs holds fully assembled endpoint URLs.
type RequestURLs struct {
	Fssp    string
	Rosstat string
	Arb     string
}

// DefaultRequestConfig returns a config populated with the previously hardcoded values.
func DefaultRequestConfig() RequestConfig {
	return RequestConfig{
		Fssp: FsspConfig{
			Endpoint: defaultFsspEndpoint,
			Inn:      "7712040126",
			Format:   "1",
			Key:      "2268a80e1f11a48f8657c69c71f1c00d41be9219",
		},
		Rosstat: RosstatConfig{
			Endpoint: defaultRosstatEndpoint,
			Inn:      "7712040126",
			Key:      "67266ba78d7779083310826cc491faf420858d1f",
		},
		Arb: ArbConfig{
			Endpoint: defaultArbEndpoint,
			Query:    "7713076301",
			Format:   "1",
			Key:      "418439f8dd7b24abbeb677bafa479280511cd9b4",
		},
	}
}

// ConfigureRequests runs the Charm form, allowing the operator to adjust every request.
func ConfigureRequests() (RequestConfig, error) {
	cfg := DefaultRequestConfig()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("ФССП").
				Description("Настройте конечный адрес и параметры запроса ФССП."),
			huh.NewInput().
				Title("FSSP endpoint").
				Value(&cfg.Fssp.Endpoint).
				Validate(required("endpoint")),
			huh.NewInput().
				Title("ИНН").
				Value(&cfg.Fssp.Inn).
				Validate(required("ИНН")),
			huh.NewInput().
				Title("format").
				Value(&cfg.Fssp.Format).
				Validate(required("format")),
			huh.NewInput().
				Title("API key").
				Value(&cfg.Fssp.Key).
				Validate(required("API key")),
		),
		huh.NewGroup(
			huh.NewNote().
				Title("Росстат").
				Description("Параметры баланса по ИНН."),
			huh.NewInput().
				Title("Rosstat endpoint").
				Value(&cfg.Rosstat.Endpoint).
				Validate(required("endpoint")),
			huh.NewInput().
				Title("ИНН").
				Value(&cfg.Rosstat.Inn).
				Validate(required("ИНН")),
			huh.NewInput().
				Title("API key").
				Value(&cfg.Rosstat.Key).
				Validate(required("API key")),
		),
		huh.NewGroup(
			huh.NewNote().
				Title("Арбитраж").
				Description("Поисковый запрос по делам и API ключ."),
			huh.NewInput().
				Title("Arb endpoint").
				Value(&cfg.Arb.Endpoint).
				Validate(required("endpoint")),
			huh.NewInput().
				Title("Поисковый запрос (q)").
				Value(&cfg.Arb.Query).
				Validate(required("q")),
			huh.NewInput().
				Title("format").
				Value(&cfg.Arb.Format).
				Validate(required("format")),
			huh.NewInput().
				Title("API key").
				Value(&cfg.Arb.Key).
				Validate(required("API key")),
		),
	)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return RequestConfig{}, ErrMenuAborted
		}
		return RequestConfig{}, err
	}

	return cfg, nil
}

// URLs assembles the three requests ready to be executed.
func (cfg RequestConfig) URLs() (RequestURLs, error) {
	fsspURL, err := cfg.FsspURL()
	if err != nil {
		return RequestURLs{}, err
	}
	rosstatURL, err := cfg.RosstatURL()
	if err != nil {
		return RequestURLs{}, err
	}
	arbURL, err := cfg.ArbURL()
	if err != nil {
		return RequestURLs{}, err
	}

	return RequestURLs{
		Fssp:    fsspURL,
		Rosstat: rosstatURL,
		Arb:     arbURL,
	}, nil
}

// FsspURL returns the final URL for the FSSP request.
func (cfg RequestConfig) FsspURL() (string, error) {
	return buildURL(cfg.Fssp.Endpoint, map[string]string{
		"inn":    cfg.Fssp.Inn,
		"format": cfg.Fssp.Format,
		"key":    cfg.Fssp.Key,
	})
}

// RosstatURL returns the final URL for the Rosstat request.
func (cfg RequestConfig) RosstatURL() (string, error) {
	return buildURL(cfg.Rosstat.Endpoint, map[string]string{
		"inn": cfg.Rosstat.Inn,
		"key": cfg.Rosstat.Key,
	})
}

// ArbURL returns the final URL for the arbitration request.
func (cfg RequestConfig) ArbURL() (string, error) {
	return buildURL(cfg.Arb.Endpoint, map[string]string{
		"q":      cfg.Arb.Query,
		"format": cfg.Arb.Format,
		"key":    cfg.Arb.Key,
	})
}

func buildURL(base string, params map[string]string) (string, error) {
	base = strings.TrimSpace(base)
	if base == "" {
		return "", fmt.Errorf("empty base url")
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("parse %q: %w", base, err)
	}

	values := parsed.Query()
	for k, v := range params {
		values.Set(k, v)
	}
	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func required(field string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("%s не может быть пустым", field)
		}
		return nil
	}
}
