package config

import (
	"github.com/spf13/viper"
)

type Collect struct {
	From string `mapstructure:"CDN_COLLECT_FROM"`
}

func LoadCollectConfig(path string, name string, typeC string) (*Collect, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(typeC)

	viper.AutomaticEnv()

	var c Collect
	c.setDefault()
	if err := viper.ReadInConfig(); err != nil {
		return &c, err
	}
	viper.Unmarshal(&c)
	return &c, nil
}

func (c *Collect) setDefault() {
	c.From = "http://localhost:4440/debug/vars"
}
