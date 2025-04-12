package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting server microservice...")

	router := gin.Default()

	router.GET("/login", func(c *gin.Context) {
		// OAuth login endpoint: redirect to the OAuth provider's authorization URL.
		c.Redirect(http.StatusTemporaryRedirect, "https://oauth.provider.com/authorize")
	})

	router.GET("/logout", func(c *gin.Context) {
		// OAuth logout endpoint: perform logout logic (e.g., clear session cookies).
		c.String(http.StatusOK, "Successfully logged out")
	})

	// Start the HTTP server on port 8080.
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
