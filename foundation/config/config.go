package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	CDNApiKey string `mapstructure:"CDN_API_KEY"`
	// StripeSecretKey string `mapstructure:"STRIPE_SECRET_KEY"`
	// PanelAddress    string `mapstructure:"PANEL_ADDRESS"`
	// Auth0Domain     string `mapstructure:"AUTH0_DOMAIN"`
	// Auth0Audience   string `mapstructure:"AUTH0_AUDIENCE"`
	// TargetUrl       string `mapstructure:"TARGET_URL"`

	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	ShutdownTimeout    time.Duration
	APIHost            string
	DebugHost          string
	CORSAllowedOrigins []string
	DefaultKID         string
	KeysFolder         string
	ActiveKID          string
	Issuer             string
	User               string
	Password           string
	HostPort           string
	Name               string
	Schema             string
	MaxIdleConns       int
	MaxOpenConns       int
	DisableTLS         bool
	ReporterURI        string
	ServiceName        string
	Probability        float64
	Build              string
	Routes             string // go build -ldflags "-X main.routes=crud"
	Desc               string
	Prefix             string
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	var c Config
	c.setDefault()
	if err := viper.ReadInConfig(); err != nil {
		return &c, err
	}
	viper.Unmarshal(&c)
	return &c, nil
}

func (c *Config) setDefault() {
	c.CDNApiKey = "apikey"
	c.ReadTimeout = time.Second * 5
	c.WriteTimeout = time.Second * 10
	c.IdleTimeout = time.Second * 120
	c.ShutdownTimeout = time.Second * 20
	c.APIHost = "0.0.0.0:3000"
	c.DebugHost = "0.0.0.0:4000"
	c.KeysFolder = "zarf/keys/"
	c.DefaultKID = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
	c.ActiveKID = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
	c.Issuer = "service project"
	c.User = "postgres"
	c.Password = "postgres"
	c.HostPort = "database-service.cdn-system.svc.cluster.local"
	c.Name = "postgres"
	c.MaxIdleConns = 2
	c.MaxOpenConns = 0
	c.DisableTLS = true
	c.ReporterURI = "tempo.cdn-system.svc.cluster.local:4317"
	c.ServiceName = "cdn-api"
	c.Probability = 0.5
	c.Build = "develop"
	c.Routes = "all"
	c.Desc = "Service Project"
	c.Prefix = "CDN"
}
