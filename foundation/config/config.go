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
	*Expvar
	*Prometheus
	*Collect
	*Publish
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string, args ...string) (*Config, error) {
	var cfg Config
	cfg.setDefault()
	for _, arg := range args {
		switch arg {
		case "web":
			w, err := LoadWebConfig(path, "web", "env")
			if err != nil {
				return nil, err
			}
			cfg.Web = w
		case "auth":
			a, err := LoadAuthConfig(path, "auth", "env")
			if err != nil {
				return nil, err
			}
			cfg.Auth = a
		case "db":
			d, err := LoadDBConfig(path, "db", "env")
			if err != nil {
				return nil, err
			}
			cfg.DB = d
		case "tempo":
			t, err := LoadTempoConfig(path, "tempo", "env")
			if err != nil {
				return nil, err
			}
			cfg.Tempo = t
		case "expvar":
			e, err := LoadExpvarConfig(path, "expvar", "env")
			if err != nil {
				return nil, err
			}
			cfg.Expvar = e
		case "prometheus":
			p, err := LoadPrometheusConfig(path, "prometheus", "env")
			if err != nil {
				return nil, err
			}
			cfg.Prometheus = p
		case "collect":
			c, err := LoadCollectConfig(path, "collect", "env")
			if err != nil {
				return nil, err
			}
			cfg.Collect = c
		case "publish":
			p, err := LoadPublishConfig(path, "publish", "env")
			if err != nil {
				return nil, err
			}
			cfg.Publish = p
		}
	}

	return &cfg, nil
}

func (cfg *Config) setDefault() {
	cfg.Build = "develop"
	cfg.Routes = "all"
	cfg.Desc = "Service Project"
}
