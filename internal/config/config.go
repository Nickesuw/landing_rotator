package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"gitlab.tubecorporate.com/platform-go/core/pkg/srv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	AppEnv string `env:"APP_ENV" yaml:"APP_ENV"`
	GeoDB  string `env:"DB_PATH"`

	ClickhouseEventCollectors []srv.Options `yaml:"clickhouse_event_collectors"`
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func (c *Config) CollectorsSRV() []srv.Options {
	var records []srv.Options
	for _, chSRV := range c.ClickhouseEventCollectors {
		records = append(records, chSRV)
	}
	return records
}

func ReadConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	err = cfg.readYaml()
	return cfg, nil
}

func (c *Config) readYaml() error {
	var filename string
	if c.IsProduction() {
		filename = "production"
	} else {
		filename = "dev"
	}
	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("./configs/%s.yaml", filename))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	return err
}
