package ui

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// URLs assembles the three requests ready to be executed.
func (cfg RequestConfig) URLs() (RequestURLs, error) {
	if !cfg.RunFssp && !cfg.RunRosstat && !cfg.RunArb {
		return RequestURLs{}, errors.New("не выбрано ни одного источника для сбора")
	}

	var urls RequestURLs

	if cfg.RunFssp {
		fsspURL, err := cfg.FsspURL()
		if err != nil {
			return RequestURLs{}, err
		}
		urls.Fssp = fsspURL
	}
	if cfg.RunRosstat {
		rosstatURL, err := cfg.RosstatURL()
		if err != nil {
			return RequestURLs{}, err
		}
		urls.Rosstat = rosstatURL
	}
	if cfg.RunArb {
		arbURL, err := cfg.ArbURL()
		if err != nil {
			return RequestURLs{}, err
		}
		urls.Arb = arbURL
	}

	return urls, nil
}

// FsspURL returns the final URL for the FSSP request.
func (cfg RequestConfig) FsspURL() (string, error) {
	inn := strings.TrimSpace(cfg.Inn)
	if inn == "" {
		return "", fmt.Errorf("ИНН не может быть пустым")
	}
	return buildURL(cfg.Fssp.Endpoint, map[string]string{
		"inn":    inn,
		"format": cfg.Fssp.Format,
		"key":    cfg.Fssp.Key,
	})
}

// RosstatURL returns the final URL for the Rosstat request.
func (cfg RequestConfig) RosstatURL() (string, error) {
	inn := strings.TrimSpace(cfg.Inn)
	if inn == "" {
		return "", fmt.Errorf("ИНН не может быть пустым")
	}
	return buildURL(cfg.Rosstat.Endpoint, map[string]string{
		"inn": inn,
		"key": cfg.Rosstat.Key,
	})
}

// ArbURL returns the final URL for the arbitration request.
func (cfg RequestConfig) ArbURL() (string, error) {
	query := strings.TrimSpace(cfg.Arb.Query)
	if query == "" {
		query = strings.TrimSpace(cfg.Inn)
	}
	if query == "" {
		return "", fmt.Errorf("поисковый запрос пустой")
	}
	return buildURL(cfg.Arb.Endpoint, map[string]string{
		"q":      query,
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
