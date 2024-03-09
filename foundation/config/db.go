package config

import (
	"github.com/spf13/viper"
)

type DB struct {
	User         string `mapstructure:"CDN_DB_USER"`
	Password     string `mapstructure:"CDN_DB_PASSWORD"`
	HostPort     string `mapstructure:"CDN_DB_HOST_PORT"`
	Name         string `mapstructure:"CDN_DB_NAME"`
	MaxIdleConns int    `mapstructure:"CDN_DB_MAX_IDLE_CONNS"`
	MaxOpenConns int    `mapstructure:"CDN_DB_MAX_OPEN_CONNS"`
	DisableTLS   bool   `mapstructure:"CDN_DB_DISABLE_TLS"`
	Schema       string `mapstructure:"CDN_DB_SCHEMA"`
}

func LoadDBConfig(path string, name string, t string) (*DB, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(t)

	viper.AutomaticEnv()

	var db DB
	db.setDefault()
	if err := viper.ReadInConfig(); err != nil {
		return &db, err
	}
	viper.Unmarshal(&db)
	return &db, nil
}

func (d *DB) setDefault() {
	d.User = "postgres"
	d.Password = "postgres"
	d.HostPort = "database-service.cdn-system.svc.cluster.local"
	d.Name = "postgres"
	d.MaxIdleConns = 2
	d.MaxOpenConns = 0
	d.DisableTLS = true
}
