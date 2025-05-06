package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/shivamkedia17/roshnii/shared/pkg/jwt"
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

const (
	// UserContextKey is the key used to store user claims in the Gin context.
	UserContextKey = "userClaims"
)

// AuthMiddleware creates a Gin middleware for JWT authentication.
func AuthMiddleware(jwtService jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Cookies: ", c.Request.Cookies())
		// 1. Get auth token from cookie
		cookie, err := c.Cookie(jwt.AuthTokenCookie)
		if err == http.ErrNoCookie {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		log.Println("Auth Token: ", cookie)
		// Validate the token
		claims, err := jwtService.ValidateToken(cookie)
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
			c.SetCookie(jwt.AuthTokenCookie, "", -1, "/", "", false, true)

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
func GetUserClaims(c *gin.Context) *jwt.Claims {
	claims, exists := c.Get(UserContextKey)
	if !exists {
		return nil
	}
	log.Printf("Claim Exist: %v", claims)

	userClaims, ok := claims.(*jwt.Claims)
	if !ok {
		return nil
	}

	log.Printf("Claim Found: %v", userClaims)
	return userClaims
}

// GetUserID retrieves the UserID from the Gin context.
// Returns empty string if claims are not found or invalid.
func GetUserID(c *gin.Context) models.UserID {
	claims := GetUserClaims(c)
	if claims == nil {
		return ""
	}
	return claims.UserID
}
