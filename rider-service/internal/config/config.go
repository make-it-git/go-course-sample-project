package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ListenAddr string `mapstructure:"LISTEN_ADDR"`
	ListenPort int    `mapstructure:"LISTEN_PORT"`
	Env        string `mapstructure:"ENV"`
}

func (c Config) ListenAddrAndPort() string {
	return fmt.Sprintf("%s:%d", c.ListenAddr, c.ListenPort)
}

func FromEnv() (*Config, error) {
	v := viper.New()
	v.SetDefault("LISTEN_ADDR", "127.0.0.1")
	v.SetDefault("LISTEN_PORT", 8000)
	v.SetDefault("ENV", "local")
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