package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shivamkedia17/roshnii/shared/pkg/models" // Adjust import path
)

// JWTService handles JWT generation and validation.
type JWTService interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// Claims defines the JWT claims structure.
type Claims struct {
	UserID models.UserID `json:"user_id"`
	Email  string        `json:"email"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTService creates a new JWT service.
func NewJWTService(secret string, duration time.Duration) JWTService {
	if secret == "" {
		panic("JWT secret cannot be empty")
	}
	return &jwtService{
		secretKey:     []byte(secret),
		tokenDuration: duration,
	}
}

// GenerateToken creates a new JWT for the given user.
func (s *jwtService) GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(s.tokenDuration)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "roshnii-service", // Optional: identify the issuer
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// ValidateToken verifies the token signature and claims.
func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Additional check: ensure UserID is present
	if claims.UserID == 0 {
		return nil, fmt.Errorf("invalid token: missing user ID")
	}

	return claims, nil
}
