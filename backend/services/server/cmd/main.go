package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors" // Import CORS middleware
	"github.com/gin-gonic/gin"

	"github.com/shivamkedia17/roshnii/services/server/internal/auth" // Adjust import paths
	"github.com/shivamkedia17/roshnii/services/server/internal/handlers"
	"github.com/shivamkedia17/roshnii/services/server/internal/middleware"
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(".") // Load from current directory (where server runs) or use "../.." for project root
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize Database Connection
	database, err := db.NewPostgresStore(cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// 3. Initialize Blob Storage (Placeholder)
	// storageService, err := storage.NewBlobStorage(cfg) // Implement this based on config.BlobStorageType
	// if err != nil {
	// 	log.Fatalf("Failed to initialize blob storage: %v", err)
	// }

	// 4. Initialize Services
	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.TokenDuration)
	googleOAuthService := auth.NewGoogleOAuthService(cfg, database, jwtService)

	// 5. Initialize Handlers
	imageHandler := handlers.NewImageHandler(database /*storageService,*/, cfg)

	// 6. Setup Gin Engine
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// 7. Setup Middleware
	// CORS - Adjust origins as needed for your frontend
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.FrontendURL, "http://localhost:8080"} // Add server origin if needed, or be more specific
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization") // Ensure Authorization header is allowed
	router.Use(cors.New(corsConfig))

	// Custom Auth Middleware
	authMW := middleware.AuthMiddleware(jwtService)

	// 8. Setup Routes
	// Base endpoint for health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	api := router.Group("/api")
	{
		// --- Authentication Routes ---
		authRoutes := api.Group("/auth") // Group under /auth
		{
			googleRoutes := authRoutes.Group("/google")
			{
				googleRoutes.GET("/login", googleOAuthService.HandleLogin)
				googleRoutes.GET("/callback", googleOAuthService.HandleCallback)
				googleRoutes.POST("/logout", authMW, googleOAuthService.HandleLogout) // Apply auth middleware here
			}

			// --- Development Only Login ---
			// WARNING: THIS IS INSECURE AND SHOULD ONLY BE USED FOR LOCAL DEV/TESTING
			if cfg.Environment == "development" {
				devRoutes := authRoutes.Group("/dev")
				{
					devRoutes.POST("/login", func(c *gin.Context) {
						var req struct {
							Email string `json:"email" binding:"required,email"`
							Name  string `json:"name"` // Optional name
						}
						if err := c.ShouldBindJSON(&req); err != nil {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
							return
						}

						// Use the new DB function
						userName := req.Name
						if userName == "" {
							userName = "Dev User" // Default name if not provided
						}
						user, err := database.FindOrCreateUserByEmail(c.Request.Context(), req.Email, userName, "dev")
						if err != nil {
							log.Printf("Dev Login Error: Failed to find/create user %s: %v", req.Email, err)
							c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login"})
							return
						}

						// Generate JWT
						token, err := jwtService.GenerateToken(user)
						if err != nil {
							log.Printf("Dev Login Error: Failed to generate JWT for user %d: %v", user.ID, err)
							c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session"})
							return
						}

						log.Printf("Dev Login Success: Generated token for user %s (ID: %d)", user.Email, user.ID)
						c.JSON(http.StatusOK, gin.H{"token": token})
					})
				}
			} else {
				log.Println("INFO: Development login endpoint is disabled in non-development environments.")
			}
		} // End /auth group

		// --- Image Routes (already includes auth middleware via RegisterRoutes) ---
		imageHandler.RegisterRoutes(api, authMW)

		// --- User Info Route (Example) ---
		api.GET("/me", authMW, func(c *gin.Context) {
			claims := middleware.GetUserClaims(c)
			if claims == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
				return
			}
			// Fetch full user details from DB if needed, or just return claims
			c.JSON(http.StatusOK, gin.H{
				"user_id": claims.UserID,
				"email":   claims.Email,
				// Add more fields if stored in claims or fetched from DB
			})
		})
	}

	// 9. Start Server
	serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Starting server on %s (Env: %s)", serverAddr, cfg.Environment)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
