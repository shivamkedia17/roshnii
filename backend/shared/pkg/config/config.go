package config

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
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

	PublicHost string `mapstructure:"PUBLIC_HOST"`
	PublicPort string `mapstructure:"PUBLIC_PORT"`

	// --- Databases ---
	PostgresURL string `mapstructure:"POSTGRES_URL"`
	QdrantURL   string `mapstructure:"QDRANT_URL"`

	// --- Storage ---
	BlobStorageType  string `mapstructure:"BLOB_STORAGE_TYPE"`
	BlobBucket       string `mapstructure:"BLOB_BUCKET"`
	AWSRegion        string `mapstructure:"AWS_REGION"`
	LocalstoragePath string `mapstructure:"LOCAL_STORAGE_PATH"`

	// Authentication
	JWTSecret        string        `mapstructure:"JWT_SECRET"`
	JWTRefreshSecret string        `mapstructure:"JWT_REFRESH_SECRET"` // Add this line
	TokenDurationStr string        `mapstructure:"TOKEN_DURATION"`
	TokenDuration    time.Duration `mapstructure:"-"`

	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	FrontendURL        string `mapstructure:"FRONTEND_URL"` // To redirect back after OAuth
}

// LoadConfig loads configuration from file and environment variables.
func LoadConfig(path string) (*Config, error) {
	// Set defaults
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SERVER_HOST", "0.0.0.0") // Default host
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("PUBLIC_HOST", "127.0.0.1") // Default host
	viper.SetDefault("PUBLIC_PORT", "8080")
	viper.SetDefault("TOKEN_DURATION", "24h")
	viper.SetDefault("BLOB_STORAGE_TYPE", "local")
	viper.SetDefault("LOCAL_STORAGE_PATH", "./uploads")
	viper.SetDefault("FRONTEND_URL", "http://localhost:5173") // Default frontend URL

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // Expecting .env or app.env file

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("")

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
		config.JWTSecret = os.Getenv("JWT_SECRET")
		log.Printf("Attempting direct env read for JWT_SECRET: found=%v", config.JWTSecret != "")
		if config.JWTSecret == "" {
			log.Fatal("JWT_SECRET must be set")
		}
	}

	// Check refresh secret
	if config.JWTRefreshSecret == "" {
		config.JWTRefreshSecret = os.Getenv("JWT_REFRESH_SECRET")
		if config.JWTRefreshSecret == "" {
			// Generate a different refresh secret if not set
			refreshBytes := make([]byte, 32)
			_, err := rand.Read(refreshBytes)
			if err != nil {
				log.Fatal("Failed to generate JWT_REFRESH_SECRET")
			}
			config.JWTRefreshSecret = base64.StdEncoding.EncodeToString(refreshBytes)
			log.Println("Warning: JWT_REFRESH_SECRET not set, using generated value")
		}
	}

	if config.GoogleClientID == "" {
		config.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
		log.Printf("Attempting direct env read for GOOGLE_CLIENT_ID: found=%v", config.GoogleClientID != "")
		if config.GoogleClientID == "" {
			log.Fatal("GOOGLE_CLIENT_ID must be set")
		}
	}

	if config.GoogleClientSecret == "" {
		config.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
		log.Printf("Attempting direct env read for GOOGLE_CLIENT_SECRET: found=%v", config.GoogleClientSecret != "")
		if config.GoogleClientSecret == "" {
			log.Fatal("GOOGLE_CLIENT_SECRET must be set")
		}
	}

	if config.PostgresURL == "" {
		config.PostgresURL = os.Getenv("POSTGRES_URL")
		log.Printf("Attempting direct env read for POSTGRES_URL: found=%v", config.PostgresURL != "")
		if config.PostgresURL == "" {
			log.Fatal("POSTGRES_URL must be set")
		}
	}

	log.Println("Configuration loaded successfully.")
	// log.Printf("Loaded config (excluding secrets): Env=%s, Port=%s", config.Environment, config.ServerPort)

	return &config, nil
}
