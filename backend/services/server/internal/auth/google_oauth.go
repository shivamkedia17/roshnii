package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/shivamkedia17/roshnii/shared/pkg/config" // Adjust import paths
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
)

const oauthStateCookieName = "oauthstate"

// GoogleOAuthService handles the Google OAuth2 flow.
type GoogleOAuthService struct {
	Config     *oauth2.Config
	UserStore  db.UserStore
	JWTService JWTService
	AppConfig  *config.Config // Store App config for FrontendURL etc.
}

// NewGoogleOAuthService initializes the Google OAuth service.
func NewGoogleOAuthService(cfg *config.Config, userStore db.UserStore, jwtService JWTService) *GoogleOAuthService {

	// redirectURL := fmt.Sprintf("http://%s:%s/api/auth/google/callback", cfg.ServerHost, cfg.ServerPort)
	redirectURL := fmt.Sprintf("http://%s:%s/api/auth/google/callback", cfg.PublicHost, cfg.PublicPort)
	log.Printf("Using Google OAuth Redirect URL: %s", redirectURL)

	return &GoogleOAuthService{
		Config: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile", // Request profile info (name, picture)
			},
			Endpoint: google.Endpoint,
		},
		UserStore:  userStore,
		JWTService: jwtService,
		AppConfig:  cfg,
	}
}

// generateState generates a random string for the OAuth state parameter.
func generateState() (string, error) {
	b := make([]byte, 32) // Increased size for better security
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes for state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// HandleLogin initiates the Google OAuth process by redirecting the user.
func (s *GoogleOAuthService) HandleLogin(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		log.Printf("Error generating OAuth state: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate login"})
		return
	}

	secureCookie := s.AppConfig.Environment != "development"
	sameSiteMode := http.SameSiteNoneMode
	if !secureCookie {
		sameSiteMode = http.SameSiteLaxMode
	}

	// Store the state in a secure, HttpOnly cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    state,
		HttpOnly: true,
		Secure:   secureCookie,                        // Use Secure in prod/staging
		Path:     "/",                                 // Accessible across the domain
		MaxAge:   int(10 * time.Minute / time.Second), // 10 minutes validity
		SameSite: sameSiteMode,                        // Good default for OAuth redirects
	})

	// Redirect user to Google's consent page
	url := s.Config.AuthCodeURL(state, oauth2.AccessTypeOffline) // Request refresh token if needed later
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleCallback handles the callback from Google after user authorization.
func (s *GoogleOAuthService) HandleCallback(c *gin.Context) {
	// 1. Verify State
	storedState, err := c.Cookie(oauthStateCookieName)
	if err != nil {
		log.Printf("OAuth Callback Error: State cookie not found: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session state. Please try logging in again."})
		return
	}
	// Clear the state cookie once used
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   oauthStateCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Delete cookie
	})

	queryState := c.Query("state")
	if queryState == "" || queryState != storedState {
		log.Printf("OAuth Callback Error: Invalid state parameter. Expected '%s', got '%s'", storedState, queryState)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OAuth state."})
		return
	}

	// 2. Handle Potential Errors from Google
	errorParam := c.Query("error")
	if errorParam != "" {
		errorDesc := c.Query("error_description")
		log.Printf("OAuth Callback Error from Google: %s - %s", errorParam, errorDesc)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed: " + errorParam})
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

	token, err := s.Config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("OAuth Callback Error: Failed to exchange code for token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login."})
		return
	}
	// Note: You might want to store token.RefreshToken securely if you need offline access later.

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
	user, err := s.UserStore.FindOrCreateUserByGoogleID(c.Request.Context(), googleUser)
	if err != nil {
		log.Printf("OAuth Callback Error: Database operation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user data."})
		return
	}

	// --- Added Logging ---
	log.Printf("------------------------------------------")
	log.Printf("OAuth Callback: Found/Created Internal User:")
	log.Printf("  User ID:       %d", user.ID)
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

	// 7. Set JWT cookies
	secureCookie := s.AppConfig.Environment != "development"

	// Set access token cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   secureCookie,
		Path:     "/",
		MaxAge:   int(s.AppConfig.TokenDuration / time.Second),
		SameSite: http.SameSiteLaxMode,
	})

	// Set refresh token cookie with longer expiration
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   secureCookie,
		Path:     "/",
		MaxAge:   int(s.AppConfig.TokenDuration * 7 / time.Second), // 7x longer
		SameSite: http.SameSiteLaxMode,
	})

	// --- Added Logging ---
	log.Printf("OAuth Callback: Set auth_token cookie successfully for user %d.", user.ID)

	// 8. Redirect to Frontend
	log.Printf("OAuth successful for user %s (%d). Redirecting to frontend: %s", user.Email, user.ID, s.AppConfig.FrontendURL)
	c.Redirect(http.StatusTemporaryRedirect, s.AppConfig.FrontendURL) // Redirect to the main page or dashboard
}

// fetchGoogleUserInfo uses the OAuth token to get user details from Google.
func (s *GoogleOAuthService) fetchGoogleUserInfo(token *oauth2.Token) (*models.GoogleUser, error) {
	client := s.Config.Client(context.Background(), token)
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

// HandleLogout clears the authentication cookies and blacklists the token.
func (s *GoogleOAuthService) HandleLogout(c *gin.Context) {
	// Get current tokens before deleting cookies
	accessToken, _ := c.Cookie("auth_token")
	refreshToken, _ := c.Cookie("refresh_token")

	// If no cookie token, try header (for dev mode)
	if accessToken == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

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

	// Clear the auth_token cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete cookie now
		HttpOnly: true,
		Secure:   s.AppConfig.Environment != "development",
		SameSite: http.SameSiteLaxMode,
	})

	// Clear the refresh_token cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete cookie now
		HttpOnly: true,
		Secure:   s.AppConfig.Environment != "development",
		SameSite: http.SameSiteLaxMode,
	})

	log.Printf("User logged out successfully.")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// HandleRefreshToken processes token refresh requests
func (s *GoogleOAuthService) HandleRefreshToken(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		// If no cookie, try getting from Authorization header (for dev/testing)
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			refreshToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not provided"})
		return
	}

	// Validate refresh token
	claims, err := s.JWTService.ValidateRefreshToken(refreshToken)
	if err != nil {
		log.Printf("Refresh token validation failed: %v", err)

		// Clear the invalid refresh token cookie
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   s.AppConfig.Environment != "development",
			SameSite: http.SameSiteLaxMode,
		})

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Get the user from database to ensure it still exists
	ctx := c.Request.Context()
	user, err := s.UserStore.GetUserByID(ctx, claims.UserID)
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
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    newAccessToken,
		HttpOnly: true,
		Secure:   s.AppConfig.Environment != "development",
		Path:     "/",
		MaxAge:   int(s.AppConfig.TokenDuration / time.Second),
		SameSite: http.SameSiteLaxMode,
	})

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"token":   newAccessToken, // Include in response for dev mode
	})
}
