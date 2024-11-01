package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database  DatabaseConfig
	WebServer WebServerConfig
	App       AppConfig
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	Schema   string
}

type WebServerConfig struct {
	Port int
}

type AppConfig struct {
	MigrationPath string
	LogLevel      string
	FilePath      string
	Env           string
}

func LoadConfig() (*Config, error) {

	v := viper.New()

	configFile := (".env")
	//configFile := filepath.Join(rootDir, "..", ".env")

	v.SetConfigFile(configFile)
	v.AutomaticEnv()

	fmt.Println("Config file:", v.ConfigFileUsed())

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No .env file found. Using default values and environment variables.")
		} else {
			fmt.Printf("error reading config file: %w\n", err)
		}
	}

	var config Config

	// Database configuration
	config.Database = DatabaseConfig{
		Driver:   getStringOrDefault(v, "DB_DRIVER", "postgres"),
		Host:     getStringOrDefault(v, "DB_HOST", "192.168.1.5"),
		Port:     getIntOrDefault(v, "DB_PORT", 5432),
		User:     getStringOrDefault(v, "DB_USERNAME", "dbuser"),
		Password: getStringOrDefault(v, "DB_PASSWORD", ""),
		Name:     getStringOrDefault(v, "DB_NAME", "mytestdb"),
		Schema:   getStringOrDefault(v, "DB_SCHEMA", "public"),
	}

	// Server configuration
	config.WebServer = WebServerConfig{
		Port: getIntOrDefault(v, "WEBSERVER_PORT", 8888),
	}

	config.App = AppConfig{
		LogLevel:      getStringOrDefault(v, "LOG_LEVEL", "debug"),
		Env:           getStringOrDefault(v, "ENV", "local"),
		MigrationPath: getStringOrDefault(v, "MIGRATION_PATH", "./migrations"),
		FilePath:      getStringOrDefault(v, "LOG_FILE_PATH", "./logs/app.log"),
	}

	return &config, nil
}

func (c *Config) GetSlogLevel() slog.Level {
	logLevel := strings.ToUpper(c.App.LogLevel)
	switch logLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo // Default log level
	}
}

func getIntOrDefault(v *viper.Viper, key string, defaultValue int) int {
	value := v.GetInt(key)
	if value == 0 && !v.IsSet(key) {
		return defaultValue
	}
	return value
}

func getStringOrDefault(v *viper.Viper, key string, defaultValue string) string {
	value := v.GetString(key)
	if value == "" && !v.IsSet(key) {
		return defaultValue
	}
	return value
}

func getBoolOrDefault(v *viper.Viper, key string, defaultValue bool) bool {
	if !v.IsSet(key) {
		return defaultValue
	}
	return v.GetBool(key)
}

func getUintOrDefault(v *viper.Viper, key string, defaultValue uint) uint {
	if !v.IsSet(key) {
		return defaultValue
	}
	return v.GetUint(key)
}
