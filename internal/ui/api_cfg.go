package ui

import (
	"fmt"
	"net/url"
	"strings"
)

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
