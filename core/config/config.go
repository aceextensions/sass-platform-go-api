package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Env              string `mapstructure:"NODE_ENV"`
	Port             string `mapstructure:"PORT"`
	DatabaseURL      string `mapstructure:"DATABASE_URL"`
	AuditDatabaseURL string `mapstructure:"AUDIT_DATABASE_URL"`
	JWTSecret        string `mapstructure:"JWT_SECRET"`
	MinioEndpoint    string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey   string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey   string `mapstructure:"MINIO_SECRET_KEY"`
}

var GlobalConfig *Config

func Load() *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values
	viper.SetDefault("NODE_ENV", "development")
	viper.SetDefault("PORT", "4000")
	viper.SetDefault("DATABASE_URL", "postgresql://aceextension:aceextension_dev@localhost:5432/aceextension")
	viper.SetDefault("AUDIT_DATABASE_URL", "postgresql://aceextension:aceextension_audit@localhost:5433/aceextension_audit")
	viper.SetDefault("JWT_SECRET", "supersecretjwtkey")

	config := &Config{}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	GlobalConfig = config
	return config
}
