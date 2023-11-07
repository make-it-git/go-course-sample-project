package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ListenAddr  string `mapstructure:"LISTEN_ADDR"`
	ListenPort  int    `mapstructure:"LISTEN_PORT"`
	Env         string `mapstructure:"ENV"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
	BrokerURL   string `mapstructure:"BROKER_URL"`
}

func (c Config) ListenAddrAndPort() string {
	return fmt.Sprintf("%s:%d", c.ListenAddr, c.ListenPort)
}

func FromEnv() (*Config, error) {
	v := viper.New()
	v.SetDefault("LISTEN_ADDR", "127.0.0.1")
	v.SetDefault("LISTEN_PORT", 8001)
	v.SetDefault("ENV", "local")
	v.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres")
	v.SetDefault("BROKER_URL", "127.0.0.1:16379")
	v.SetConfigType("env")
	v.AutomaticEnv()

	cfg := Config{}
	err := v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	if _, ok := allEnvs[cfg.Env]; !ok {
		return nil, NewUnknownEnvErr(cfg.Env)
	}

	return &cfg, nil
}
