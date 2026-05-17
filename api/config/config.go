package config

import "os"

var (
	ApiPort               = os.Getenv("API_PORT")
	ApiVersion            = os.Getenv("API_VERSION")
	AuthSecret            = os.Getenv("AUTH_SECRET")
	Env                   = environment(os.Getenv("ENV"))
	GoTrueURL             = os.Getenv("GOTRUE_URL")
	GoTrueServiceKey      = os.Getenv("GOTRUE_SERVICE_KEY")
	StripeSecretKey       = os.Getenv("STRIPE_SECRET_KEY")
	StripePublishableKey  = os.Getenv("STRIPE_PUBLISHABLE_KEY")
	StripeWebhookSecret   = os.Getenv("STRIPE_WEBHOOK_SECRET")
	ProductionDatabase = &DatabaseConfig{
		Host:         os.Getenv("DATABASE_HOST"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
		User:         os.Getenv("DATABASE_USER"),
		Password:     os.Getenv("DATABASE_PASSWORD"),
		Port:         os.Getenv("DATABASE_PORT"),
	}
)

type DatabaseConfig struct {
	Host string
	User string
	// Password contains database password - gosec ignore
	Password     string // #nosec G117
	DatabaseName string
	Port         string
}

func (d DatabaseConfig) String() string {
	return "host=" + d.Host + " user=" + d.User + " port=" + d.Port + " database=" + d.DatabaseName + " password=**"
}
