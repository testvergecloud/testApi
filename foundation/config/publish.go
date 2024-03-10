package config

import (
	"time"

	"github.com/spf13/viper"
)

type Publish struct {
	To       string        `mapstructure:"CDN_PUBLISH_TO"`
	Interval time.Duration `mapstructure:"CDN_PUBLISH_INTERVAL"`
}

func LoadPublishConfig(path string, name string, t string) (*Publish, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(t)

	viper.AutomaticEnv()

	var p Publish
	p.setDefault()
	if err := viper.ReadInConfig(); err != nil {
		return &p, err
	}
	viper.Unmarshal(&p)
	return &p, nil
}

func (p *Publish) setDefault() {
	p.To = "console"
	p.Interval = time.Second * 5
}
