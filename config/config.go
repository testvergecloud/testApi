package config

import "github.com/spf13/viper"

type Config struct {
	CDNApiKey       string `mapstructure:"CDN_API_KEY"`
	StripeSecretKey string `mapstructure:"STRIPE_SECRET_KEY"`
	PanelAddress    string `mapstructure:"PANEL_ADDRESS"`
	Auth0Domain     string `mapstructure:"AUTH0_DOMAIN"`
	Auth0Audience   string `mapstructure:"AUTH0_AUDIENCE"`
	TargetUrl       string `mapstructure:"TARGET_URL"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	var c Config
	if err := viper.ReadInConfig(); err != nil {
		return &c, err
	}
	viper.Unmarshal(&c)
	return &c, nil
}
