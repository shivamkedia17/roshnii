package routes

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shivamkedia17/roshnii/services/server/internal/handlers"
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
)

func RegisterAuthRoutes(routerGroup *gin.RouterGroup, authMiddleware gin.HandlerFunc, h *handlers.GoogleOAuthService) {
	googleRoutes := routerGroup.Group("/auth/google")
	{
		googleRoutes.GET("/login", h.HandleLogin)
		googleRoutes.GET("/callback", h.HandleCallback)
		googleRoutes.POST("/refresh", h.HandleRefreshToken)
		googleRoutes.POST("/logout", authMiddleware, h.HandleLogout)
	}
}

// RegisterRoutes connects album routes to the Gin engine.
func RegisterAlbumRoutes(routerGroup *gin.RouterGroup, authMiddleware gin.HandlerFunc, h *handlers.AlbumHandler) {
	albumRoutes := routerGroup.Group("/albums")
	albumRoutes.Use(authMiddleware)
	{
		albumRoutes.POST("", h.CreateAlbum)
		albumRoutes.GET("", h.ListAlbums)
		albumRoutes.GET("/:id", h.GetAlbum)
		albumRoutes.PUT("/:id", h.UpdateAlbum)
		albumRoutes.DELETE("/:id", h.DeleteAlbum)
		albumRoutes.GET("/:id/images", h.ListAlbumImages)
		albumRoutes.POST("/:id/images", h.AddImageToAlbum)
		albumRoutes.DELETE("/:id/images/:image_id", h.RemoveImageFromAlbum)
	}
}

func RegisterImageRoutes(routerGroup *gin.RouterGroup, authMiddleware gin.HandlerFunc, h *handlers.ImageHandler) {
	imageRoutes := routerGroup.Group("/images")
	imageRoutes.Use(authMiddleware)
	{
		imageRoutes.GET("", h.HandleListImages)                 // List user images
		imageRoutes.POST("/upload", h.HandleUploadImage)        // Upload endpoint
		imageRoutes.GET("/:id", h.HandleGetImage)               // Single image metadata
		imageRoutes.DELETE("/:id", h.HandleDeleteImage)         // Delete image
		imageRoutes.GET("/:id/download", h.HandleDownloadImage) // Download image file
	}
}

func RegisterUserRoutes(routerGroup *gin.RouterGroup, authMiddleware gin.HandlerFunc, h *handlers.UserHandler) {
	userRoutes := routerGroup.Group("/me")
	userRoutes.Use(authMiddleware)
	{
		userRoutes.GET("", authMiddleware, h.GetCurrentUser) // User Profile Info Endpoint
	}

	// If you want to add more user-related endpoints:
	// router.PUT("/me", authMiddleware, h.UpdateUserProfile)
	// router.GET("/users/:id", authMiddleware, h.GetUserByID) // For public profiles, if needed
}

// TODO
// func RegisterSearchRoutes(routerGroup *gin.RouterGroup, authMiddleware gin.HandlerFunc, h *handlers.SearchHandler) {
// 	searchRoutes := routerGroup.Group("/search")
// 	searchRoutes.Use(authMiddleware)
// 	{
// 		searchRoutes.GET("", h.Search)
// 	}
// }

func SetupRouter(cfg *config.Config, handlers *handlers.Handlers, authMiddleware gin.HandlerFunc) *gin.Engine {
	// a. Extract Relevant Config
	environment := cfg.Environment
	frontEndURL := cfg.FrontendURL
	frontendBuildPath := cfg.FrontendBuildPath
	// TODO // serverURL 	:= deps.Config.ServerHost

	// b. Setup Gin Engine
	if environment == config.ProdEnvironment {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// c. Setup CORS (Cross Origin Resource Sharing)
	corsConfig := cors.DefaultConfig()

	// Set allowed origins from configuration
	// FIXME change hardcoded server address to cfg based dynamic address
	corsConfig.AllowOrigins = []string{frontEndURL, "http://127.0.0.1:8080", "http://127.0.0.1:5173", "http://localhost:5173", "http://localhost:8080"}

	// Enable credentials for cookies
	corsConfig.AllowCredentials = true

	// // Allow common methods
	// corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	// // Allow common headers
	// corsConfig.AllowHeaders = []string{
	// 	"Origin", "Content-Type", "Accept", "Authorization",
	// 	"X-Requested-With", "X-CSRF-Token", "Access-Control-Allow-Origin",
	// }

	// Add Access-Control-Expose-Headers to expose custom headers to the frontend
	corsConfig.ExposeHeaders = []string{"Content-Length", "Content-Type"}

	router.Use(cors.New(corsConfig))

	// d. Setup Routes

	// Base endpoint for health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// API routes
	api := router.Group("/api")

	RegisterAuthRoutes(api, authMiddleware, &handlers.OAuth)
	RegisterImageRoutes(api, authMiddleware, &handlers.Img)
	RegisterAlbumRoutes(api, authMiddleware, &handlers.Album)
	RegisterUserRoutes(api, authMiddleware, &handlers.User)
	// RegisterSearchRoutes()

	// Serve frontend static files
	RegisterStaticAssets(router, frontendBuildPath)

	return router
}

// RegisterStaticAssets sets up routes to serve the React frontend static files
func RegisterStaticAssets(router *gin.Engine, frontendPath string) {
	// Serve the static files (JS, CSS, images)
	router.Static("/assets", frontendPath+"/assets")

	// Serve the favicon and other root files
	// router.StaticFile("/favicon.ico", frontendPath+"/favicon.ico")
	// router.StaticFile("/robots.txt", frontendPath+"/robots.txt")

	// For any other route, serve the index.html file (for React router to handle)
	router.NoRoute(func(c *gin.Context) {
		// Only serve index.html for non-API routes
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.File(frontendPath + "/index.html")
		} else {
			// Let API 404s pass through
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
		}
	})
}
