package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Environment defines the running environment (e.g., "development", "staging", "production")
	Environment string `mapstructure:"ENVIRONMENT"`

	LogLevel string `mapstructure:"LOG_LEVEL"`

	ServerHost string `mapstructure:"SERVER_HOST"` // Host for callback URL etc.
	ServerPort string `mapstructure:"SERVER_PORT"`

	// --- Databases ---
	PostgresURL string `mapstructure:"POSTGRES_URL"`
	QdrantURL   string `mapstructure:"QDRANT_URL"`

	// --- Storage ---
	BlobStorageType  string `mapstructure:"BLOB_STORAGE_TYPE"`
	BlobBucket       string `mapstructure:"BLOB_BUCKET"`
	AWSRegion        string `mapstructure:"AWS_REGION"`
	LocalstoragePath string `mapstructure:"LOCAL_STORAGE_PATH"`

	// --- Authentication ---
	JWTSecret        string        `mapstructure:"JWT_SECRET"`
	TokenDurationStr string        `mapstructure:"TOKEN_DURATION"` // Use string for viper
	TokenDuration    time.Duration `mapstructure:"-"`              // Calculated duration

	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	FrontendURL        string `mapstructure:"FRONTEND_URL"` // To redirect back after OAuth
}

// LoadConfig loads configuration from file and environment variables.
func LoadConfig(path string) (*Config, error) {
	// Set defaults
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SERVER_HOST", "localhost") // Default host
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("TOKEN_DURATION", "24h")
	viper.SetDefault("BLOB_STORAGE_TYPE", "local")
	viper.SetDefault("LOCAL_STORAGE_PATH", "./uploads")
	viper.SetDefault("FRONTEND_URL", "http://localhost:3000") // Default frontend URL

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // Expecting .env or app.env file

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Println("Config file not found, using environment variables and defaults.")
	} else if err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Printf("Unable to decode config into struct: %v", err)
		return nil, err
	}

	// Calculate TokenDuration from string
	duration, err := time.ParseDuration(config.TokenDurationStr)
	if err != nil {
		log.Printf("Invalid TOKEN_DURATION format: %v. Using default 24h.", err)
		duration = 24 * time.Hour
	}
	config.TokenDuration = duration

	// Basic validation
	if config.JWTSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	if config.GoogleClientID == "" || config.GoogleClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set")
	}
	if config.PostgresURL == "" {
		log.Fatal("POSTGRES_URL must be set")
	}

	log.Println("Configuration loaded successfully.")
	// log.Printf("Loaded config (excluding secrets): Env=%s, Port=%s", config.Environment, config.ServerPort)

	return &config, nil
}
