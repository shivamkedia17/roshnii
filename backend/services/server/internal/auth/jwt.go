package auth

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shivamkedia17/roshnii/shared/pkg/models" // Adjust import path
)

// JWTService handles JWT generation and validation.
type JWTService interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	GenerateRefreshToken(user *models.User) (string, error)
	ValidateRefreshToken(tokenString string) (*Claims, error)
	BlacklistToken(tokenString string) error
	IsTokenAboutToExpire(claims *Claims, thresholdMinutes int) bool
	RefreshToken(claims *Claims) (string, error) // To refresh the short lived token using the refresh token
}

// Claims defines the JWT claims structure.
type Claims struct {
	UserID     models.UserID `json:"user_id"`
	Email      string        `json:"email"`
	Name       string        `json:"name,omitempty"`
	PictureURL string        `json:"picture_url,omitempty"`
	TokenType  string        `json:"token_type,omitempty"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey         []byte
	refreshSecretKey  []byte
	tokenDuration     time.Duration
	refreshDuration   time.Duration
	blacklistedTokens map[string]time.Time // Simple in-memory blacklist
	blacklistMutex    sync.Mutex           // Add a mutex to protect the map
}

// NewJWTService creates a new JWT service.
func NewJWTService(secret string, refreshSecret string, duration time.Duration) JWTService {
	if secret == "" {
		panic("JWT secret cannot be empty")
	}

	// Refresh token lasts longer (e.g., 7 days)
	refreshDuration := duration * 7

	return &jwtService{
		secretKey:         []byte(secret),
		refreshSecretKey:  []byte(refreshSecret),
		tokenDuration:     duration,
		refreshDuration:   refreshDuration,
		blacklistedTokens: make(map[string]time.Time),
	}
}

// GenerateToken creates a new JWT for the given user.
func (s *jwtService) GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(s.tokenDuration)

	// Enhanced claims with more user information
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		PictureURL: func() string {
			if user.PictureURL != nil {
				return *user.PictureURL
			} else {
				return ""
			}
		}(),
		TokenType: "access", // Specify this is an access token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "roshnii-service",
			Subject:   fmt.Sprintf("%s", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// GenerateRefreshToken creates a refresh token with longer expiration
func (s *jwtService) GenerateRefreshToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(s.refreshDuration)
	claims := &Claims{
		UserID:    user.ID,
		Email:     user.Email,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "roshnii-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.refreshSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}
	return tokenString, nil
}

// ValidateToken verifies the token signature and claims.
func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	// Check if token is blacklisted
	if s.IsTokenBlacklisted(tokenString) {
		return nil, fmt.Errorf("token is blacklisted")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	}, jwt.WithValidMethods([]string{"HS256"})) // Explicitly specify valid signing methods

	if err != nil {
		// Enhanced error handling for expired tokens
		if err.Error() == "token has invalid claims: token is expired" {
			return nil, fmt.Errorf("token expired: %w", err)
		} else if strings.Contains(err.Error(), "used before issued") {
			return nil, fmt.Errorf("token not yet valid: %w", err)
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Verify it's an access token, not a refresh token
	if claims.TokenType != "" && claims.TokenType != "access" {
		return nil, fmt.Errorf("invalid token type: expected access token")
	}

	// Enhanced validation
	// Additional check: ensure UserID is present
	if claims.UserID == "" {
		return nil, fmt.Errorf("invalid token: missing user ID")
	}

	// Check if email is present
	if claims.Email == "" {
		return nil, fmt.Errorf("invalid token: missing email")
	}

	return claims, nil
}

func (s *jwtService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	// Check if token is blacklisted
	if s.IsTokenBlacklisted(tokenString) {
		return nil, fmt.Errorf("refresh token is blacklisted")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.refreshSecretKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Verify it's a refresh token
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type: expected refresh token")
	}

	if claims.UserID == "" {
		return nil, fmt.Errorf("invalid token: missing user ID")
	}

	return claims, nil
}

// Implement BlacklistToken method
func (s *jwtService) BlacklistToken(tokenString string) error {
	// Parse token without validation to get expiration time
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return fmt.Errorf("failed to parse token for blacklisting: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok {
		if claims.ExpiresAt != nil {
			// Lock before modifying the map
			s.blacklistMutex.Lock()
			// Add to blacklist with expiry time
			s.blacklistedTokens[tokenString] = claims.ExpiresAt.Time
			// Clean up expired tokens from blacklist (optional)
			s.cleanupBlacklist()
			s.blacklistMutex.Unlock()
			return nil
		}
	}

	// If no expiry found, blacklist for a default period (token duration)
	s.blacklistMutex.Lock()
	s.blacklistedTokens[tokenString] = time.Now().Add(s.tokenDuration)
	s.blacklistMutex.Unlock()
	return nil
}

// Implement IsTokenBlacklisted helper
func (s *jwtService) IsTokenBlacklisted(tokenString string) bool {
	s.blacklistMutex.Lock()
	defer s.blacklistMutex.Unlock()

	expiry, exists := s.blacklistedTokens[tokenString]

	// If token exists in blacklist and hasn't expired in the blacklist
	if exists && time.Now().Before(expiry) {
		return true
	}

	// If it exists but has expired, remove it from blacklist
	if exists {
		delete(s.blacklistedTokens, tokenString)
	}

	return false
}

// Implement cleanupBlacklist helper
func (s *jwtService) cleanupBlacklist() {
	// Assume caller already hols the lock
	now := time.Now()
	for token, expiry := range s.blacklistedTokens {
		if now.After(expiry) {
			delete(s.blacklistedTokens, token)
		}
	}
}

// IsTokenAboutToExpire checks if a token is about to expire within the given threshold
func (s *jwtService) IsTokenAboutToExpire(claims *Claims, thresholdMinutes int) bool {
	if claims == nil || claims.ExpiresAt == nil {
		return true
	}

	expiryTime := claims.ExpiresAt.Time
	thresholdDuration := time.Duration(thresholdMinutes) * time.Minute

	// Token is about to expire if the time remaining is less than the threshold
	return time.Until(expiryTime) < thresholdDuration
}

// RefreshToken creates a new token with the same claims but a new expiry time
func (s *jwtService) RefreshToken(claims *Claims) (string, error) {
	if claims == nil {
		return "", fmt.Errorf("cannot refresh nil claims")
	}

	// Create a new claims object with refreshed timestamps
	newClaims := &Claims{
		UserID:     claims.UserID,
		Email:      claims.Email,
		Name:       claims.Name,
		PictureURL: claims.PictureURL,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    claims.Issuer,
			Subject:   claims.Subject,
		},
	}

	// Create a new token with refreshed claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign refreshed token: %w", err)
	}

	return tokenString, nil
}
