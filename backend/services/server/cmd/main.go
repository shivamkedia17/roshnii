package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors" // Import CORS middleware
	"github.com/gin-gonic/gin"

	"github.com/shivamkedia17/roshnii/services/server/internal/auth" // Adjust import paths
	"github.com/shivamkedia17/roshnii/services/server/internal/handlers"
	"github.com/shivamkedia17/roshnii/services/server/internal/middleware"
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"

	"github.com/shivamkedia17/roshnii/shared/pkg/storage" // Add this import
)

func main() {
	// 1. Load Configuration
	// Load from current directory (where server runs) or use "../.." for project root
	cfg, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize Database Connection
	database, err := db.NewPostgresStore(cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// 3. Initialize Blob Storage
	var storageService storage.BlobStorage
	if cfg.BlobStorageType == "local" {
		localStoragePath := cfg.LocalstoragePath
		if localStoragePath == "" {
			localStoragePath = "./uploads" // Default path
		}

		var err error
		storageService, err = storage.NewLocalStorage(localStoragePath)
		if err != nil {
			log.Fatalf("Failed to initialize local storage: %v", err)
		}
		log.Printf("Using local file storage at: %s", localStoragePath)
	} else {
		// For now, fall back to local storage if type is unrecognized
		log.Printf("Unrecognized storage type '%s', using local storage", cfg.BlobStorageType)
		storageService, _ = storage.NewLocalStorage("./uploads")
	}

	// 4. Initialize Services
	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTRefreshSecret, cfg.TokenDuration)
	googleOAuthService := auth.NewGoogleOAuthService(cfg, database, jwtService)

	// 5. Initialize Handlers
	imageHandler := handlers.NewImageHandler(database, storageService, cfg)
	albumHandler := handlers.NewAlbumHandler(database, cfg)

	// TODO IMP
	// searchHandler := handlers.NewSearchHandler(database, cfg)

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

	// Ensure Authorization header is allowed
	// corsConfig.AddAllowHeaders("Authorization")

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
		// a. Authentication Routes
		authRoutes := api.Group("/auth")
		{
			googleRoutes := authRoutes.Group("/google")
			{
				googleRoutes.GET("/login", googleOAuthService.HandleLogin)
				googleRoutes.GET("/callback", googleOAuthService.HandleCallback)
				googleRoutes.POST("/logout", authMW, googleOAuthService.HandleLogout) // Apply auth middleware here
				googleRoutes.POST("/refresh", googleOAuthService.HandleRefreshToken)
			}

			// --- Development Only Login ---
			// This section is compiled only in development mode and provides a simplified
			// login flow for local development and testing purposes
			if cfg.Environment == "development" {
				// Create a dev routes group with a clear warning
				devRoutes := authRoutes.Group("/dev")

				// Add a middleware to log warnings about dev-only endpoints
				devRoutes.Use(func(c *gin.Context) {
					log.Println("WARNING: Development-only authentication endpoint accessed")
					c.Next()
				})

				// Dev login endpoint - NEVER USE IN PRODUCTION
				devRoutes.POST("/login", func(c *gin.Context) {
					// Check if we're really in development mode as a safety measure
					if cfg.Environment != "development" {
						c.JSON(http.StatusForbidden, gin.H{"error": "Development login routes are disabled in production"})
						return
					}

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

					// Generate JWT tokens
					accessToken, err := jwtService.GenerateToken(user)
					if err != nil {
						log.Printf("Dev Login Error: Failed to generate access token for user %s: %v", user.ID, err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session"})
						return
					}

					refreshToken, err := jwtService.GenerateRefreshToken(user)
					if err != nil {
						log.Printf("Dev Login Error: Failed to generate refresh token for user %s: %v", user.ID, err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session"})
						return
					}

					// Set cookies - Use same approach as the OAuth flow
					cookieSecure := cfg.Environment != "development"

					http.SetCookie(c.Writer, &http.Cookie{
						Name:     "auth_token",
						Value:    accessToken,
						HttpOnly: true,
						Secure:   cookieSecure,
						Path:     "/",
						MaxAge:   int(cfg.TokenDuration / time.Second),
						SameSite: http.SameSiteLaxMode,
					})

					http.SetCookie(c.Writer, &http.Cookie{
						Name:     "refresh_token",
						Value:    refreshToken,
						HttpOnly: true,
						Secure:   cookieSecure,
						Path:     "/",
						MaxAge:   int(cfg.TokenDuration * 7 / time.Second),
						SameSite: http.SameSiteLaxMode,
					})

					log.Printf("Dev Login Success: Generated tokens for user %s (ID: %s)", user.Email, user.ID)

					// Return success response with minimal information
					// Only include token information in response for development testing
					c.JSON(http.StatusOK, gin.H{
						"message":    "Development login successful",
						"user_id":    user.ID,
						"email":      user.Email,
						"expires_in": int(cfg.TokenDuration / time.Second),
						// The tokens are in cookies, but also return them for dev testing
						"dev_note":      "TOKENS PROVIDED FOR DEVELOPMENT TESTING ONLY",
						"token":         accessToken,
						"refresh_token": refreshToken,
					})
				})
			}
		}

		// End /auth group

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

		// --- Image Routes (already includes auth middleware via RegisterRoutes) ---
		imageHandler.RegisterRoutes(api, authMW)
		albumHandler.RegisterRoutes(api, authMW)

		// TODO
		// searchHandler.RegisterRoutes(api, authMW)

	}

	// 9. Start Server
	serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Starting server on %s (Env: %s)", serverAddr, cfg.Environment)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
