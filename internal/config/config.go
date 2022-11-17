package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type environment string

const (
	prod environment = "prod"
	dev              = "dev"
)

func (e environment) IsProd() bool {
	return e == prod
}

func (e environment) IsDev() bool {
	return e == dev
}

func NewConfig() (*AppConfig, error) {
	c := &AppConfig{}
	err := envconfig.Process("RFLAT", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

type AppConfig struct {
	WorkerCount  int           `default:"10" split_words:"true"`
	PollInterval time.Duration `default:"1m" split_words:"true"`
	Environment  environment   `default:"prod"`
}
