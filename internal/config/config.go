package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Host        string `mapstructure:"HOST" doc:"The hostname for the server."`
	Port        int    `mapstructure:"PORT" doc:"The port number for the server."`
	Env         string `mapstructure:"ENV" doc:"The environment in which the app is running (e.g., dev, staging, prod)."`
	FrontendURL string `mapstructure:"FRONTEND_URL" doc:"The URL of the frontend application."`

	DB struct {
		DSN          string `mapstructure:"DB_BLOG_DSN" doc:"The PostgreSQL connection string."`
		MaxOpenConns int    `mapstructure:"DB_MAX_OPEN_CONNS" doc:"The maximum number of open connections to the database."`
		MaxIdleConns int    `mapstructure:"DB_MAX_IDLE_CONNS" doc:"The maximum number of idle connections to the database."`
		MaxIdleTime  string `mapstructure:"DB_MAX_IDLE_TIME" doc:"The maximum duration a connection can remain idle (e.g., '5m')."`
	} `doc:"Database configuration."`

	Limiter struct {
		RPS     float64 `mapstructure:"LIMITER_RPS" doc:"The rate limit in requests per second."`
		Burst   int     `mapstructure:"LIMITER_BURST" doc:"The maximum number of burst requests allowed."`
		Enabled bool    `mapstructure:"LIMITER_ENABLED" doc:"Whether the rate limiter is enabled."`
	} `doc:"Rate limiter configuration."`

	SMTP struct {
		Host     string `mapstructure:"SMTP_HOST" doc:"The SMTP server hostname."`
		Port     int    `mapstructure:"SMTP_PORT" doc:"The SMTP server port."`
		Username string `mapstructure:"SMTP_USERNAME" doc:"The SMTP username for authentication."`
		Password string `mapstructure:"SMTP_PASSWORD" doc:"The SMTP password for authentication."`
		Sender   string `mapstructure:"SMTP_SENDER" doc:"The email address of the sender."`
	} `doc:"SMTP email configuration."`

	CORS struct {
		TrustedOrigins []string `mapstructure:"CORS_TRUSTED_ORIGINS" doc:"Space-separated list of trusted origins for CORS."`
	} `doc:"CORS configuration."`

	Upload struct {
		Path     string `mapstructure:"UPLOAD_PATH" doc:"The directory path for permanent file uploads."`
		TempPath string `mapstructure:"UPLOAD_TEMP_PATH" doc:"The directory path for temporary file uploads."`
	} `doc:"File upload configuration."`
}

func SetConfigDefaultValues() {
	viper.SetDefault("HOST", "localhost")
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("ENV", "development")
	viper.SetDefault("FRONTEND_URL", "http://localhost:3000")

	viper.SetDefault("DB_BLOG_DSN", "user:password@tcp(localhost:3306)/blogdb")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_TIME", "15m")

	viper.SetDefault("LIMITER_RPS", 2)
	viper.SetDefault("LIMITER_BURST", 4)
	viper.SetDefault("LIMITER_ENABLED", true)

	viper.SetDefault("SMTP_HOST", "sandbox.smtp.mailtrap.io")
	viper.SetDefault("SMTP_PORT", 25)
	viper.SetDefault("SMTP_USERNAME", "your_username")
	viper.SetDefault("SMTP_PASSWORD", "your_password")
	viper.SetDefault("SMTP_SENDER", "Example <noreply@example.com>")

	viper.SetDefault("CORS_TRUSTED_ORIGINS", "http://localhost:3000 http://localhost:8080")
	viper.SetDefault("UPLOAD_PATH", "./upload")
	viper.SetDefault("UPLOAD_TEMP_PATH", "./temp")
}

func LoadConfig(cfg *Config) error {
	// Map Viper keys to config struct fields
	cfg.Host = viper.GetString("HOST")
	cfg.Port = viper.GetInt("PORT")
	cfg.Env = viper.GetString("ENV")
	cfg.FrontendURL = viper.GetString("FRONTEND_URL")

	cfg.DB.DSN = viper.GetString("DB_BLOG_DSN")
	cfg.DB.MaxOpenConns = viper.GetInt("DB_MAX_OPEN_CONNS")
	cfg.DB.MaxIdleConns = viper.GetInt("DB_MAX_IDLE_CONNS")
	cfg.DB.MaxIdleTime = viper.GetString("DB_MAX_IDLE_TIME")

	cfg.Limiter.RPS = viper.GetFloat64("LIMITER_RPS")
	cfg.Limiter.Burst = viper.GetInt("LIMITER_BURST")
	cfg.Limiter.Enabled = viper.GetBool("LIMITER_ENABLED")

	cfg.SMTP.Host = viper.GetString("SMTP_HOST")
	cfg.SMTP.Port = viper.GetInt("SMTP_PORT")
	cfg.SMTP.Username = viper.GetString("SMTP_USERNAME")
	cfg.SMTP.Password = viper.GetString("SMTP_PASSWORD")
	cfg.SMTP.Sender = viper.GetString("SMTP_SENDER")

	cfg.Upload.Path = viper.GetString("UPLOAD_PATH")
	cfg.Upload.TempPath = viper.GetString("UPLOAD_TEMP_PATH")

	// Trusted origins env is space-separated string; convert to []string
	trustedOrigins := viper.GetString("CORS_TRUSTED_ORIGINS")
	cfg.CORS.TrustedOrigins = strings.Fields(trustedOrigins)

	return nil
}
