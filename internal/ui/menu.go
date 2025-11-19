package ui

import (
	"errors"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

const (
	defaultFsspEndpoint    = "https://api.damia.ru/fssp/isps"
	defaultRosstatEndpoint = "https://api.damia.ru/rs/balance"
	defaultArbEndpoint     = "https://api.damia.ru/arb/dela"
)

var ErrMenuAborted = errors.New("request configuration canceled by user")

type serviceChoice string

const (
	serviceFssp    serviceChoice = "fssp"
	serviceRosstat serviceChoice = "rosstat"
	serviceArb     serviceChoice = "arb"
)

type RequestConfig struct {
	Inn        string
	RunFssp    bool
	RunRosstat bool
	RunArb     bool

	Fssp    FsspConfig
	Rosstat RosstatConfig
	Arb     ArbConfig
}

// FsspConfig describes the FSSP API request.
type FsspConfig struct {
	Endpoint string
	Format   string
	Key      string
}

// RosstatConfig describes the Rosstat API request.
type RosstatConfig struct {
	Endpoint string
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
	defaultInn := "inn"
	return RequestConfig{
		Inn:        defaultInn,
		RunFssp:    true,
		RunRosstat: true,
		RunArb:     true,
		Fssp: FsspConfig{
			Endpoint: defaultFsspEndpoint,
			Format:   "1",
			Key:      "2268a80e1f11a48f8657c69c71f1c00d41be9219",
		},
		Rosstat: RosstatConfig{
			Endpoint: defaultRosstatEndpoint,
			Key:      "67266ba78d7779083310826cc491faf420858d1f",
		},
		Arb: ArbConfig{
			Endpoint: defaultArbEndpoint,
			Query:    defaultInn,
			Format:   "1",
			Key:      "418439f8dd7b24abbeb677bafa479280511cd9b4",
		},
	}
}

// ConfigureRequests runs the Charm form, allowing the operator to adjust every request.
func ConfigureRequests() (RequestConfig, error) {
	cfg := DefaultRequestConfig()
	selectedServices := defaultServiceSelection(cfg)
	advancedMode := false

	formKeyMap := huh.NewDefaultKeyMap()
	formKeyMap.Confirm.Toggle = key.NewBinding(
		key.WithKeys("ctrl+a", "h", "l", "left", "right"),
		key.WithHelp("ctrl+a", "показать/скрыть расширенные"),
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Сбор данных").
				Description("Укажите один ИНН компании и выберите, какие данные собирать."),
			huh.NewInput().
				Title("ИНН предприятия").
				Value(&cfg.Inn).
				Validate(required("ИНН")),
			huh.NewMultiSelect[serviceChoice]().
				Title("Что нужно собрать").
				Description("Пробел — отметить, Enter — продолжить").
				Options(
					huh.NewOption("ФССП", serviceFssp),
					huh.NewOption("Росстат", serviceRosstat),
					huh.NewOption("Арбитраж", serviceArb),
				).
				Value(&selectedServices).
				Validate(func(values []serviceChoice) error {
					if len(values) == 0 {
						return errors.New("выберите хотя бы один источник данных")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Расширенные настройки (Ctrl+A)").
				Description("Нужно, чтобы изменить endpoint и ключи под каждую функцию.").
				Value(&advancedMode),
		),
		huh.NewGroup(
			huh.NewNote().
				Title("ФССП").
				Description("Настройте конечный адрес и параметры запроса ФССП."),
			huh.NewInput().
				Title("FSSP endpoint").
				Value(&cfg.Fssp.Endpoint).
				Validate(required("endpoint")),
			huh.NewInput().
				Title("format").
				Value(&cfg.Fssp.Format).
				Validate(required("format")),
			huh.NewInput().
				Title("API key").
				Value(&cfg.Fssp.Key).
				Validate(required("API key")),
		).WithHideFunc(func() bool { return !advancedMode }),
		huh.NewGroup(
			huh.NewNote().
				Title("Росстат").
				Description("Параметры баланса по ИНН."),
			huh.NewInput().
				Title("Rosstat endpoint").
				Value(&cfg.Rosstat.Endpoint).
				Validate(required("endpoint")),
			huh.NewInput().
				Title("API key").
				Value(&cfg.Rosstat.Key).
				Validate(required("API key")),
		).WithHideFunc(func() bool { return !advancedMode }),
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
				Value(&cfg.Arb.Query),
			huh.NewInput().
				Title("format").
				Value(&cfg.Arb.Format).
				Validate(required("format")),
			huh.NewInput().
				Title("API key").
				Value(&cfg.Arb.Key).
				Validate(required("API key")),
		).WithHideFunc(func() bool { return !advancedMode }),
	).WithKeyMap(formKeyMap)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return RequestConfig{}, ErrMenuAborted
		}
		return RequestConfig{}, err
	}

	cfg.RunFssp = serviceSelected(selectedServices, serviceFssp)
	cfg.RunRosstat = serviceSelected(selectedServices, serviceRosstat)
	cfg.RunArb = serviceSelected(selectedServices, serviceArb)

	return cfg, nil
}

func defaultServiceSelection(cfg RequestConfig) []serviceChoice {
	var services []serviceChoice
	if cfg.RunFssp {
		services = append(services, serviceFssp)
	}
	if cfg.RunRosstat {
		services = append(services, serviceRosstat)
	}
	if cfg.RunArb {
		services = append(services, serviceArb)
	}
	return services
}

func serviceSelected(choices []serviceChoice, target serviceChoice) bool {
	for _, choice := range choices {
		if choice == target {
			return true
		}
	}
	return false
}
