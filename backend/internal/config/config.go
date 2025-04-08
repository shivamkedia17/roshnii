package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration parameters for the application.
// It uses 'mapstructure' tags to specify the corresponding key in the config source (env vars, file).
type Config struct {
	// --- General ---
	// Environment defines the running environment (e.g., "development", "staging", "production")
	Environment string `mapstructure:"ENVIRONMENT"`

	// LogLevel sets the logging level (e.g., "debug", "info", "warn", "error")
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// --- Server (Common for multiple services) ---
	// ServerPort is the port the specific service will listen on. Each service might override this.
	ServerPort string `mapstructure:"SERVER_PORT"`

	// --- Databases ---
	PostgresURL string `mapstructure:"POSTGRES_URL"`
	QdrantURL   string `mapstructure:"QDRANT_URL"`
	// QdrantAPIKey is optional, depending on your Qdrant setup.
	QdrantAPIKey string `mapstructure:"QDRANT_API_KEY"`

	// --- Storage ---
	// BlobStorageType specifies the provider (e.g., "s3", "gcs", "azure", "local").
	BlobStorageType string `mapstructure:"BLOB_STORAGE_TYPE"`
	BlobBucket      string `mapstructure:"BLOB_BUCKET"`
	// Provider-specific settings (add more as needed)
	AWSRegion        string `mapstructure:"AWS_REGION"`         // Example for S3
	LocalstoragePath string `mapstructure:"LOCAL_STORAGE_PATH"` // Example for local dev

	// --- Authentication ---
	JWTSecret string `mapstructure:"JWT_SECRET"`
	// TokenDuration specifies how long JWT tokens are valid for (e.g., "1h", "24h").
	TokenDuration string `mapstructure:"TOKEN_DURATION"`

	// --- Service Discovery / Communication (Optional) ---
	// URLs for direct communication if not using service discovery or message queues.
	// AuthServiceURL         string `mapstructure:"AUTH_SERVICE_URL"`
	// DeduplicationServiceURL string `mapstructure:"DEDUPLICATION_SERVICE_URL"`
	// EmbeddingServiceURL    string `mapstructure:"EMBEDDING_SERVICE_URL"` // Python service
	// FaceRecServiceURL      string `mapstructure:"FACE_REC_SERVICE_URL"`    // Python service
	// SearchServiceURL       string `mapstructure:"SEARCH_SERVICE_URL"`

	// --- Message Queue (Optional) ---
	// MessageQueueURL string `mapstructure:"MESSAGE_QUEUE_URL"`
	// Add specific queue/topic names if needed
}

// LoadConfig loads configuration from file and environment variables.
// Environment variables override file settings.
// It looks for a file named '.env' or 'config.env' in the specified path.
func LoadConfig(path string) (config *Config, err error) {
	// Set defaults (optional, but good practice)
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SERVER_PORT", "8080") // Default port, can be overridden per service
	viper.SetDefault("TOKEN_DURATION", "24h")
	viper.SetDefault("BLOB_STORAGE_TYPE", "local")      // Default to local for easier dev setup
	viper.SetDefault("LOCAL_STORAGE_PATH", "./uploads") // Default local path

	// Tell viper the path/s to look for the config file in.
	viper.AddConfigPath(path) // e.g., "." for the current directory

	// Tell viper the name of the config file (without extension).
	viper.SetConfigName("app") // Will look for "app.env" or "app.yaml" etc.
	// OR specific name like ".env" - viper needs SetConfigFile explicitly then
	// viper.SetConfigFile(".env")

	// Tell viper the type of the config file.
	viper.SetConfigType("env") // ".env" file format

	// Enable reading environment variables automatically.
	viper.AutomaticEnv()
	// Set a replacer for env vars (e.g., SERVER_PORT maps to ServerPort)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Attempt to read the config file.
	err = viper.ReadInConfig()
	// Ignore error if config file is not found, env vars might be sufficient.
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Println("Config file not found, using environment variables and defaults.")
		err = nil // Reset error, not finding file is ok if env vars are set
	} else if err != nil {
		// Some other error occurred reading the config file
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	// Unmarshal the configuration into the Config struct.
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Printf("Unable to decode config into struct: %v", err)
		return nil, err
	}

	log.Println("Configuration loaded successfully.")
	// Optionally print loaded config for debugging (mask secrets!)
	// log.Printf("Loaded config: %+v", config) // Be careful with secrets!

	return config, nil
}
