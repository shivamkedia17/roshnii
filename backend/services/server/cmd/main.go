package main

import (
	"fmt"
	"log"

	// Import CORS middleware

	"github.com/shivamkedia17/roshnii/services/server/internal/handlers"
	"github.com/shivamkedia17/roshnii/services/server/internal/middleware"
	"github.com/shivamkedia17/roshnii/services/server/internal/routes"
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
	"github.com/shivamkedia17/roshnii/shared/pkg/jwt"
	"github.com/shivamkedia17/roshnii/shared/pkg/storage"
)

func main() {
	// 1. Load Configuration
	// Load from current directory (where server runs) or use "../.." for project root
	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize Database Connection
	db, err := db.NewPostgresStore(cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 3. Initialize Blob Storage
	storageService, err := storage.InitStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialise Blob Store: %v", err)
	}

	// 4. Initialize JWT Service
	jwtService := jwt.NewJWTService(cfg.JWTSecret, cfg.JWTRefreshSecret, cfg.TokenDuration)

	// 5. Initialize Handlers & Middleware
	handlers := handlers.InitHandlers(cfg, db, storageService, jwtService)
	authMiddleware := middleware.AuthMiddleware(jwtService)

	// 6. Setup Routing
	router := routes.SetupRouter(cfg, &handlers, authMiddleware)

	// 7. Start Server
	serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Starting server on %s (Env: %s)", serverAddr, cfg.Environment)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
