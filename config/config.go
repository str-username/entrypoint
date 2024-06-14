package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Entrypoint struct {
		Host string `yaml:"host" validate:"required"`
	}
	MongoDB struct {
		Url  string `yaml:"url" validate:"required"`
		Db   string `yaml:"db" validate:"required"`
		Coll string `yaml:"coll" validate:"required"`
	}
}

func (c Config) Validate() error {
	return validator.New().Struct(c)
}

type Valid interface {
	Validate() error
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	configFile, err := os.Open(configPath)

	if err != nil {
		log.Fatal().Err(err).Str("func", "NewConfig").Msg("Failed to open config file")
		return nil, err
	}

	defer configFile.Close()

	configNewDecoder := yaml.NewDecoder(configFile)

	if err = configNewDecoder.Decode(config); err != nil {
		log.Fatal().Err(err).Str("func", "NewConfig").Msg("Decode config error")
		return nil, err
	}

	var validConfig Valid = config

	if err = validConfig.Validate(); err != nil {
		log.Fatal().Err(err).Str("func", "NewConfig").Msg("Validate config error")
		return nil, err
	}

	return config, nil
}
