package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type ModuleConfig struct {
	Enabled  bool                   `mapstructure:"enabled"`
	Settings map[string]interface{} `mapstructure:"settings"`
}

type PluginConfig struct {
	Name     string                 `mapstructure:"name"`
	Enabled  bool                   `mapstructure:"enabled"`
	Priority int                    `mapstructure:"priority"`
	Settings map[string]interface{} `mapstructure:"settings"`
}

type TenancyConfig struct {
	Isolation struct {
		Mode string `mapstructure:"mode"` // shared | isolated | hybrid
	} `mapstructure:"isolation"`
	Workspaces struct {
		EnableSeparateSchemas   bool `mapstructure:"enable_separate_schemas"`
		EnableSeparateDatabases bool `mapstructure:"enable_separate_databases"`
		MultiStoreCompatibility    struct {
			Enabled          bool   `mapstructure:"enabled"`
			StoreGroupMap    string `mapstructure:"store_group_maps_to"`
			StoreViewMap     string `mapstructure:"store_view_maps_to"`
		} `mapstructure:"multi_store_compatibility"`
	} `mapstructure:"workspaces"`
}

type Config struct {
	Env              string `mapstructure:"NODE_ENV"`
	Port             string `mapstructure:"PORT"`
	DatabaseURL      string `mapstructure:"DATABASE_URL"`
	AuditDatabaseURL string `mapstructure:"AUDIT_DATABASE_URL"`
	JWTSecret        string `mapstructure:"JWT_SECRET"`
	SessionSecret    string `mapstructure:"SESSION_SECRET"`
	ApiBaseURL       string `mapstructure:"API_BASE_URL"`
	AppBaseURL       string `mapstructure:"APP_BASE_URL"`

	// OAuth
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GithubClientID     string `mapstructure:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `mapstructure:"GITHUB_CLIENT_SECRET"`

	// MinIO / Storage
	MinioEndpoint  string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	MinioUseSSL    bool   `mapstructure:"MINIO_USE_SSL"`

	// Email / SMTP
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPass     string `mapstructure:"SMTP_PASS"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`
	RequirePhone bool   `mapstructure:"REQUIRE_PHONE"`

	// Slack
	SlackWebhookURL string `mapstructure:"SLACK_WEBHOOK_URL"`

	// Feature Flags (Legacy support)
	EnableAccounting bool `mapstructure:"ENABLE_ACCOUNTING"`
	EnableInventory  bool `mapstructure:"ENABLE_INVENTORY"`
	EnableSales      bool `mapstructure:"ENABLE_SALES"`
	EnablePurchase   bool `mapstructure:"ENABLE_PURCHASE"`
	EnableCRM        bool `mapstructure:"ENABLE_CRM"`

	// Framework levels
	Modules map[string]ModuleConfig `mapstructure:"modules"`
	Plugins struct {
		Official  []PluginConfig `mapstructure:"official"`
		Community []PluginConfig `mapstructure:"community"`
	} `mapstructure:"plugins"`

	Tenancy TenancyConfig `mapstructure:"tenancy"`
}

var GlobalConfig *Config

func Load() *Config {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 1. Set Defaults
	v.SetDefault("NODE_ENV", "development")
	v.SetDefault("PORT", "4000")

	// 2. Load .env (Secrets)
	v.SetConfigFile(".env")
	if err := v.ReadInConfig(); err != nil {
		log.Printf("Note: No .env file found: %v", err)
	}

	// 3. Load config.yaml (Main Framework Config)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config") // Fallback to config directory
	if err := v.MergeInConfig(); err != nil {
		log.Printf("Note: No config.yaml file found: %v", err)
	}

	// 4. Load tenancy.yaml
	v.SetConfigName("tenancy")
	if err := v.MergeInConfig(); err != nil {
		log.Printf("Note: No tenancy.yaml file found: %v", err)
	}

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	GlobalConfig = config
	return config
}
