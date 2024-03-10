package config

import (
	"time"

	"github.com/spf13/viper"
)

type Prometheus struct {
	Host            string        `mapstructure:"CDN_PROMETHEUS_HOST"`
	Route           string        `mapstructure:"CDN_PROMETHEUS_ROUTE"`
	ReadTimeout     time.Duration `mapstructure:"CDN_PROMETHEUS_READ_TIMEOUT"`
	WriteTimeout    time.Duration `mapstructure:"CDN_PROMETHEUS_WRITE_TIMEOUT"`
	IdleTimeout     time.Duration `mapstructure:"CDN_PROMETHEUS_IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"CDN_PROMETHEUS_SHUTDOWN_TIMEOUT"`
}

func LoadPrometheusConfig(path string, name string, typeC string) (*Prometheus, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(typeC)

	viper.AutomaticEnv()

	var p Prometheus
	p.setDefault()
	if err := viper.ReadInConfig(); err != nil {
		return &p, err
	}
	viper.Unmarshal(&p)
	return &p, nil
}

func (p *Prometheus) setDefault() {
	p.Host = "0.0.0.0:3002"
	p.Route = "/metrics"
	p.ReadTimeout = time.Second * 5
	p.WriteTimeout = time.Second * 10
	p.IdleTimeout = time.Second * 120
	p.ShutdownTimeout = time.Second * 5
}
