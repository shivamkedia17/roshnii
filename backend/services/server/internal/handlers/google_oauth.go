package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/shivamkedia17/roshnii/shared/pkg/config" // Adjust import paths
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
	"github.com/shivamkedia17/roshnii/shared/pkg/jwt"
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

// generateState generates a random string for the OAuth state parameter.
func generateState() (string, error) {
	b := make([]byte, 32) // Increased size for better security
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes for state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GoogleOAuthService handles the Google OAuth2 flow.
type GoogleOAuthService struct {
	Config      *config.Config // Store App config for FrontendURL etc.
	DB          db.UserStore
	JWTService  jwt.JWTService
	OAuthConfig *oauth2.Config
}

// NewGoogleOAuthService initializes the Google OAuth service.
func NewGoogleOAuthService(cfg *config.Config, userStore db.UserStore, jwtService jwt.JWTService) *GoogleOAuthService {

	// redirectURL := fmt.Sprintf("http://%s:%s/api/auth/google/callback", cfg.ServerHost, cfg.ServerPort)

	//  don't redirect to the backend?
	redirectURL := fmt.Sprintf("http://%s:%s/api/auth/google/callback", cfg.PublicHost, cfg.PublicPort)

	// FIXME hardcoded frontend address
	// redirectURL := "http://localhost:5173"

	// CORRECT: redirect to the frontend instead
	log.Printf("Using Google OAuth Redirect URL: %s", redirectURL)

	return &GoogleOAuthService{
		OAuthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile", // Request profile info (name, picture)
			},
			Endpoint: google.Endpoint,
		},
		DB:         userStore,
		JWTService: jwtService,
		Config:     cfg,
	}
}

// fetchGoogleUserInfo uses the OAuth token to get user details from Google.
func (s *GoogleOAuthService) fetchGoogleUserInfo(token *oauth2.Token) (*models.GoogleUser, error) {
	client := s.OAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo") // v3 endpoint
	if err != nil {
		return nil, fmt.Errorf("failed to request user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google API returned non-200 status: %s - %s", resp.Status, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response body: %w", err)
	}

	var googleUser models.GoogleUser
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info JSON: %w", err)
	}

	// Validate required fields
	if googleUser.ID == "" {
		return nil, fmt.Errorf("google user info missing ID")
	}
	if googleUser.Email == "" {
		return nil, fmt.Errorf("google user info missing email")
	}
	if !googleUser.VerifiedEmail {
		return nil, fmt.Errorf("google email not verified")
	}

	return &googleUser, nil
}

// HandleLogin initiates the Google OAuth process by redirecting the user.
func (s *GoogleOAuthService) HandleLogin(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		log.Printf("Error generating OAuth state: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate login"})
		return
	}

	isCookieSecure := s.Config.Environment == config.ProdEnvironment
	maxAge := int(5 * time.Minute / time.Second) // 5 minutes validity for the state

	// Store the state in a secure, HttpOnly cookie
	c.SetCookie(jwt.StateCookie, state, maxAge, "/", s.Config.CookieDomain, isCookieSecure, true)
	c.SetSameSite(http.SameSiteLaxMode)

	// Refresh Token is NOT fetched from ID provider. Ask the user to reauthenticate when expires.
	authConsentURL := s.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)

	// WARNING: Don't Redirect user to Google's consent page
	// c.Redirect(http.StatusTemporaryRedirect, authConsentURL)

	// Instead of redirecting, return the URL to the frontend
	c.JSON(http.StatusOK, gin.H{
		"auth_url": authConsentURL,
	})
}

// HandleCallback handles the callback from Google after user authorization.
func (s *GoogleOAuthService) HandleCallback(c *gin.Context) {
	// // 1. Verify State
	// storedState, err := c.Cookie(jwt.StateCookie)
	// if err != nil {
	cookies := c.Request.Cookies()
	log.Printf("Cookies Found: %v", cookies)
	// 	log.Printf("OAuth Callback Error: State cookie not found: %v", err)
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid callback state when logging in. Please try logging in again."})
	// 	return
	// }

	// // Clear the state cookie once used
	// c.SetCookie(jwt.StateCookie, "", -1, "/", "", false, true)
	// c.SetSameSite(http.SameSiteLaxMode)

	// queryState := c.Query("state")
	// if queryState == "" || queryState != storedState {
	// 	log.Printf("OAuth Callback Error: Invalid state parameter. Expected '%s', got '%s'", storedState, queryState)
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state."})
	// 	return
	// }

	// 2. Handle Potential Errors from Google
	errorParam := c.Query("error")
	if errorParam != "" {
		errorDesc := c.Query("error_description")
		log.Printf("OAuth Callback Error from Google: %s - %s", errorParam, errorDesc)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization failed: " + errorParam})
		// Optionally redirect to a frontend error page
		// c.Redirect(http.StatusTemporaryRedirect, s.AppConfig.FrontendURL+"/login?error="+errorParam)
		return
	}

	// 3. Exchange Code for Token
	code := c.Query("code")
	if code == "" {
		log.Println("OAuth Callback Error: No code parameter in the request.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code missing."})
		return
	}

	token, err := s.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("OAuth Callback Error: Failed to exchange code for token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login."})
		return
	}

	googleUser, err := s.fetchGoogleUserInfo(token)
	if err != nil {
		log.Printf("OAuth Callback Error: Failed to get user info from Google: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information."})
		return
	}

	// --- Added Logging ---
	log.Printf("------------------------------------------")
	log.Printf("OAuth Callback: Fetched Google User Info:")
	log.Printf("  ID (sub):      %s", googleUser.ID)
	log.Printf("  Email:         %s", googleUser.Email)
	log.Printf("  Verified:      %t", googleUser.VerifiedEmail)
	log.Printf("  Name:          %s", googleUser.Name)
	log.Printf("  Picture URL:   %s", googleUser.Picture)
	log.Printf("------------------------------------------")
	// --- End Added Logging ---

	// 5. Find or Create User in DB
	user, err := s.DB.FindOrCreateUserByGoogleID(c.Request.Context(), googleUser)
	if err != nil {
		log.Printf("OAuth Callback Error: Database operation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user data."})
		return
	}

	// --- Added Logging ---
	log.Printf("------------------------------------------")
	log.Printf("OAuth Callback: Found/Created Internal User:")
	log.Printf("  User ID:       %s", user.ID)
	log.Printf("  Email:         %s", user.Email)
	log.Printf("  Name:          %s", user.Name)
	log.Printf("  Auth Provider: %s", user.AuthProvider)
	log.Printf("  Created At:    %s", user.CreatedAt.String()) // Convert time to string
	log.Printf("------------------------------------------")
	// --- End Added Logging ---

	// 6. Generate JWT tokens (access + refresh)
	accessToken, err := s.JWTService.GenerateToken(user)
	if err != nil {
		log.Printf("OAuth Callback Error: Failed to generate access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete login."})
		return
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(user)
	if err != nil {
		log.Printf("OAuth Callback Error: Failed to generate refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete login."})
		return
	}

	isCookieSecure := s.Config.Environment == config.ProdEnvironment

	// Set access token cookie
	c.SetCookie(jwt.AuthTokenCookie, accessToken, int(s.Config.TokenDuration/time.Second), "/", s.Config.CookieDomain, isCookieSecure, true)
	c.SetSameSite(http.SameSiteLaxMode)

	// Set refresh token cookie with longer expiration
	c.SetCookie(jwt.RefreshTokenCookie, refreshToken, int(s.Config.TokenDuration*7/time.Second), "/", s.Config.CookieDomain, isCookieSecure, true)
	c.SetSameSite(http.SameSiteLaxMode)

	// --- Added Logging ---
	log.Printf("OAuth Callback: Set auth_token cookie successfully for user %s.", user.ID)

	// 8. Redirect to Frontend
	log.Printf("OAuth successful for user %s (%s). Redirecting to frontend: %s", user.Email, user.ID, s.Config.FrontendURL)
	c.Redirect(http.StatusTemporaryRedirect, s.Config.FrontendURL) // Redirect to the main page or dashboard
}

// HandleLogout clears the authentication cookies and blacklists the token.
func (s *GoogleOAuthService) HandleLogout(c *gin.Context) {
	// Get current tokens from cookies
	accessToken, _ := c.Cookie(jwt.AuthTokenCookie)
	refreshToken, _ := c.Cookie(jwt.RefreshTokenCookie)

	// Blacklist tokens if they exist
	if accessToken != "" {
		if err := s.JWTService.BlacklistToken(accessToken); err != nil {
			log.Printf("Warning: Failed to blacklist access token: %v", err)
		}
	}

	if refreshToken != "" {
		if err := s.JWTService.BlacklistToken(refreshToken); err != nil {
			log.Printf("Warning: Failed to blacklist refresh token: %v", err)
		}
	}

	isCookieSecure := s.Config.Environment == config.ProdEnvironment

	// Clear the auth_token cookie
	c.SetCookie(jwt.AuthTokenCookie, "", -1, "/", s.Config.CookieDomain, isCookieSecure, true)
	c.SetSameSite(http.SameSiteLaxMode)

	// Clear the refresh_token cookie
	c.SetCookie(jwt.RefreshTokenCookie, "", -1, "/", s.Config.CookieDomain, isCookieSecure, true)
	c.SetSameSite(http.SameSiteLaxMode)

	// Add explicit cache control headers
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Header("Pragma", "no-cache")

	log.Printf("User logged out successfully.")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// HandleRefreshToken processes token refresh requests
// only accepts refresh tokens from cookies.
func (s *GoogleOAuthService) HandleRefreshToken(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie(jwt.RefreshTokenCookie)

	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not provided"})
		return
	}

	isCookieSecure := s.Config.Environment == config.ProdEnvironment

	// Validate refresh token
	claims, err := s.JWTService.ValidateRefreshToken(refreshToken)
	if err != nil {
		log.Printf("Refresh token validation failed: %v", err)

		c.SetCookie(jwt.RefreshTokenCookie, "", -1, "/", s.Config.CookieDomain, isCookieSecure, true)
		c.SetSameSite(http.SameSiteLaxMode)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Get the user from database to ensure it still exists
	ctx := c.Request.Context()
	user, err := s.DB.GetUserByID(ctx, claims.UserID)
	if err != nil {
		log.Printf("Error fetching user during token refresh: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Generate new access token
	newAccessToken, err := s.JWTService.GenerateToken(user)
	if err != nil {
		log.Printf("Failed to generate new access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Set the new access token cookie
	c.SetCookie(jwt.AuthTokenCookie, newAccessToken, int(s.Config.TokenDuration/time.Second), "/", s.Config.CookieDomain, isCookieSecure, true)
	c.SetSameSite(http.SameSiteLaxMode)

	// Response with success message
	response := gin.H{"message": "Token refreshed successfully"}

	c.JSON(http.StatusOK, response)
}
