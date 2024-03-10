package config

import (
	"time"

	"github.com/spf13/viper"
)

type Expvar struct {
	Host            string        `mapstructure:"CDN_EXPVAR_HOST"`
	Route           string        `mapstructure:"CDN_EXPVAR_ROUTE"`
	ReadTimeout     time.Duration `mapstructure:"CDN_EXPVAR_READ_TIMEOUT"`
	WriteTimeout    time.Duration `mapstructure:"CDN_EXPVAR_WRITE_TIMEOUT"`
	IdleTimeout     time.Duration `mapstructure:"CDN_EXPVAR_IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"CDN_EXPVAR_SHUTDOWN_TIMEOUT"`
}

func LoadExpvarConfig(path string, name string, typeC string) (*Expvar, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(typeC)

	viper.AutomaticEnv()

	var e Expvar
	e.setDefault()
	if err := viper.ReadInConfig(); err != nil {
		return &e, err
	}
	viper.Unmarshal(&e)
	return &e, nil
}

func (e *Expvar) setDefault() {
	e.Host = "0.0.0.0:3001"
	e.Route = "/metrics"
	e.ReadTimeout = time.Second * 5
	e.WriteTimeout = time.Second * 10
	e.IdleTimeout = time.Second * 120
	e.ShutdownTimeout = time.Second * 5
}
