package config

import (
	"github.com/spf13/viper"
)

type Tempo struct {
	ReporterURI string  `mapstructure:"CDN_TEMPO_REPORTER_URI"`
	ServiceName string  `mapstructure:"CDN_TEMPO_SERVICE_NAME"`
	Probability float64 `mapstructure:"CDN_TEMPO_PROBABILITY"`
	// Shouldn't use a high Probability value in non-developer systems.
	// 0.05 should be enough for most systems. Some might want to have
	// this even lower.
}

func LoadTempoConfig(path string, name string, typeC string) (*Tempo, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(typeC)

	viper.AutomaticEnv()

	var t Tempo
	t.setDefault()
	if err := viper.ReadInConfig(); err != nil {
		return &t, err
	}
	viper.Unmarshal(&t)
	return &t, nil
}

func (t *Tempo) setDefault() {
	t.ReporterURI = "tempo.cdn-system.svc.cluster.local:4317"
	t.ServiceName = "cdn-api"
	t.Probability = 0.05
}
