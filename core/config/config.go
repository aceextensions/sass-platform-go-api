package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Env                string `mapstructure:"NODE_ENV"`
	Port               string `mapstructure:"PORT"`
	DatabaseURL        string `mapstructure:"DATABASE_URL"`
	AuditDatabaseURL   string `mapstructure:"AUDIT_DATABASE_URL"`
	JWTSecret          string `mapstructure:"JWT_SECRET"`
	MinioEndpoint      string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey     string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey     string `mapstructure:"MINIO_SECRET_KEY"`
	AppBaseURL         string `mapstructure:"APP_BASE_URL"`
	ApiBaseURL         string `mapstructure:"API_BASE_URL"`
	RequirePhone       bool   `mapstructure:"REQUIRE_PHONE"`
	SMTPHost           string `mapstructure:"SMTP_HOST"`
	SMTPPort           int    `mapstructure:"SMTP_PORT"`
	SMTPUser           string `mapstructure:"SMTP_USER"`
	SMTPPass           string `mapstructure:"SMTP_PASS"`
	SMTPFrom           string `mapstructure:"SMTP_FROM"`
	SlackWebhookURL    string `mapstructure:"SLACK_WEBHOOK_URL"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GithubClientID     string `mapstructure:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `mapstructure:"GITHUB_CLIENT_SECRET"`
	SessionSecret      string `mapstructure:"SESSION_SECRET"`
	EnableInventory    bool   `mapstructure:"ENABLE_INVENTORY"`
	EnableAccounting   bool   `mapstructure:"ENABLE_ACCOUNTING"`
	EnableSales        bool   `mapstructure:"ENABLE_SALES"`
	EnablePurchase     bool   `mapstructure:"ENABLE_PURCHASE"`
	EnableCRM          bool   `mapstructure:"ENABLE_CRM"`
}

var GlobalConfig *Config

func Load() *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicitly read .env file if it exists
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read .env file: %v", err)
	}

	// Default values
	viper.SetDefault("NODE_ENV", "development")
	viper.SetDefault("PORT", "4000")
	viper.SetDefault("DATABASE_URL", "postgresql://aceextension:aceextension_dev@localhost:5432/aceextension")
	viper.SetDefault("AUDIT_DATABASE_URL", "postgresql://aceextension:aceextension_audit@localhost:5433/aceextension_audit")
	viper.SetDefault("JWT_SECRET", "supersecretjwtkey")
	viper.SetDefault("ENABLE_INVENTORY", false)
	viper.SetDefault("ENABLE_ACCOUNTING", false)
	viper.SetDefault("ENABLE_SALES", false)
	viper.SetDefault("ENABLE_PURCHASE", false)
	viper.SetDefault("ENABLE_CRM", false)

	config := &Config{}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	GlobalConfig = config
	return config
}
