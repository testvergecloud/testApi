package config

import (
	"github.com/spf13/viper"
)

type Auth struct {
	KeysFolder string `mapstructure:"CDN_AUTH_KEYS_FOLDER"`
	ActiveKID  string `mapstructure:"CDN_AUTH_ACTIVE_KID"`
	DefaultKID string `mapstructure:"CDN_AUTH_DEFAULT_KID"`
	Issuer     string `mapstructure:"CDN_AUTH_ISSUER"`
}

func LoadAuthConfig(path string, name string, typeC string) (*Auth, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(typeC)

	viper.AutomaticEnv()

	var a Auth
	a.setDefault()
	if err := viper.ReadInConfig(); err != nil {
		return &a, err
	}
	viper.Unmarshal(&a)
	return &a, nil
}

func (a *Auth) setDefault() {
	a.KeysFolder = "zarf/keys/"
	a.DefaultKID = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
	a.ActiveKID = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
	a.Issuer = "service project"
}
