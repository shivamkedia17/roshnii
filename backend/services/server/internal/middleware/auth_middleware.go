package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/shivamkedia17/roshnii/services/server/internal/auth" // Adjust import paths
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

const (
	// UserContextKey is the key used to store user claims in the Gin context.
	UserContextKey = "userClaims"
)

// AuthMiddleware creates a Gin middleware for JWT authentication.
func AuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get token from Cookie (preferred) or Authorization header
		tokenString := ""
		cookie, err := c.Cookie("auth_token")
		if err == nil && cookie != "" {
			tokenString = cookie
		} else {
			// Fallback to Authorization header (e.g., "Bearer <token>")
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
				return
			}
			tokenString = parts[1]
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication token not found"})
			return
		}

		// Validate the token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			errorMessage := "Invalid or expired token"
			statusCode := http.StatusUnauthorized

			// More specific error messages based on the error
			if strings.Contains(err.Error(), "token is blacklisted") {
				errorMessage = "Session has been invalidated, please log in again"
			} else if strings.Contains(err.Error(), "expired") {
				errorMessage = "Session expired, please refresh your token or log in again"
				statusCode = http.StatusUnauthorized
			} else if strings.Contains(err.Error(), "invalid token") {
				errorMessage = "Invalid authentication token"
			}

			log.Printf("Token validation failed: %v", err)

			// Clear potentially invalid cookie
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "auth_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})

			c.AbortWithStatusJSON(statusCode, gin.H{"error": errorMessage})
			return
		}

		// 3. Store claims in context for handlers to access
		c.Set(UserContextKey, claims)

		// Continue processing the request
		c.Next()
	}
}

// GetUserClaims retrieves user claims from the Gin context.
// Returns nil if claims are not found or invalid.
func GetUserClaims(c *gin.Context) *auth.Claims {
	claims, exists := c.Get(UserContextKey)
	if !exists {
		return nil
	}
	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		return nil
	}
	return userClaims
}

// GetUserID retrieves the UserID from the Gin context.
// Returns 0 if claims are not found or invalid.
func GetUserID(c *gin.Context) models.UserID {
	claims := GetUserClaims(c)
	if claims == nil {
		return 0
	}
	return claims.UserID
}
