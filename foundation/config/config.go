package config

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	CDNApiKey string `mapstructure:"CDN_API_KEY"`
	// StripeSecretKey string `mapstructure:"STRIPE_SECRET_KEY"`
	// PanelAddress    string `mapstructure:"PANEL_ADDRESS"`
	// Auth0Domain     string `mapstructure:"AUTH0_DOMAIN"`
	// Auth0Audience   string `mapstructure:"AUTH0_AUDIENCE"`
	// TargetUrl       string `mapstructure:"TARGET_URL"`
	Build  string
	Routes string // go build -ldflags "-X main.routes=crud"
	Desc   string
	*Web
	*Auth
	*DB
	*Tempo
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (*Config, error) {
	w, err := LoadWebConfig(path, "web", "env")
	if err != nil {
		return nil, err
	}
	a, err := LoadAuthConfig(path, "auth", "env")
	if err != nil {
		return nil, err
	}
	db, err := LoadDBConfig(path, "db", "env")
	if err != nil {
		return nil, err
	}
	t, err := LoadTempoConfig(path, "tempo", "env")
	if err != nil {
		return nil, err
	}

	c := Config{
		Web:   w,
		Auth:  a,
		DB:    db,
		Tempo: t,
	}
	c.setDefault()

	return &c, nil
}

func (c *Config) setDefault() {
	c.Build = "develop"
	c.Routes = "all"
	c.Desc = "Service Project"
}
