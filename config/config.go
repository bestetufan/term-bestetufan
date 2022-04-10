package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GINMode        string        `mapstructure:"GIN_MODE"`
	LogLevel       string        `mapstructure:"LOG_LEVEL"`
	AutoMigrate    bool          `mapstructure:"AUTO_MIGRATE"`
	HTTPPort       string        `mapstructure:"HTTP_PORT"`
	DBHost         string        `mapstructure:"DB_HOST"`
	DBPort         string        `mapstructure:"DB_PORT"`
	DBDatabaseName string        `mapstructure:"DB_NAME"`
	DBUsername     string        `mapstructure:"DB_USERNAME"`
	DBPassword     string        `mapstructure:"DB_PASSWORD"`
	JWTSecret      string        `mapstructure:"JWT_SECRET"`
	JWTIss         string        `mapstructure:"JWT_ISS"`
	JWTExp         time.Duration `mapstructure:"JWT_EXP"`
}

// Load reads configuration from environment variables.
func Load() (config *Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
