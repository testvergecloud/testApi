package config

import (
	"time"

	"github.com/spf13/viper"
)

type Web struct {
	ReadTimeout        time.Duration `mapstructure:"CDN_WEB_READ_TIMEOUT"`
	WriteTimeout       time.Duration `mapstructure:"CDN_WEB_WRITE_TIMEOUT"`
	IdleTimeout        time.Duration `mapstructure:"CDN_WEB_IDLE_TIMEOUT"`
	ShutdownTimeout    time.Duration `mapstructure:"CDN_WEB_SHUTDOWN_TIMEOUT"`
	APIHost            string        `mapstructure:"CDN_WEB_API_HOST"`
	DebugHost          string        `mapstructure:"CDN_WEB_DEBUG_HOST"`
	CORSAllowedOrigins []string      `mapstructure:"CDN_WEB_CORS_ALLOWED_ORIGIN"`
}

func LoadWebConfig(path string, name string, typeC string) (*Web, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(typeC)

	viper.AutomaticEnv()

	var w Web
	w.setDefault()
	if err := viper.ReadInConfig(); err != nil {
		return &w, err
	}
	viper.Unmarshal(&w)
	return &w, nil
}

func (w *Web) setDefault() {
	w.ReadTimeout = time.Second * 5
	w.WriteTimeout = time.Second * 10
	w.IdleTimeout = time.Second * 120
	w.ShutdownTimeout = time.Second * 20
	w.APIHost = "0.0.0.0:3000"
	w.DebugHost = "0.0.0.0:4000"
	w.CORSAllowedOrigins = []string{"*"}
}
