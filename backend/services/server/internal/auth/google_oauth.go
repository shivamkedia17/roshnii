package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"honnef.co/go/tools/config"
)

// GoogleOAuthService holds the OAuth2 configuration along with dependencies.
type GoogleOAuthService struct {
	Config     *oauth2.Config
	UserStore  UserStore
	JWTService JWTService
}

// UserStore represents user database methods to find or create a user.
type UserStore interface {
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

// JWTService represents the JWT generation service.
type JWTService interface {
	GenerateToken(userID int64, email string) (string, error)
}

// NewGoogleOAuthService initializes the Google OAuth service.
func NewGoogleOAuthService(cfg *config.Config, userStore UserStore, jwtService JWTService) *GoogleOAuthService {
	return &GoogleOAuthService{
		Config: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  fmt.Sprintf("http://localhost:%s/api/auth/google/callback", cfg.ServerPort),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
		UserStore:  userStore,
		JWTService: jwtService,
	}
}

// generateState generates a random string for the OAuth state parameter.
func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// OAuthLogin initiates the Google OAuth process.
func (s *GoogleOAuthService) OAuthLogin(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	// Typically store the state in session or cookie for later verification.
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
	})
	url := s.Config.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleUser represents the user information from Google.
type GoogleUser struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// OAuthCallback handles the callback from Google.
func (s *GoogleOAuthService) OAuthCallback(c *gin.Context) {
	// Verify state parameter using the cookie.
	storedState, err := c.Cookie("oauthstate")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State cookie not found"})
		return
	}
	state := c.Query("state")
	if state != storedState {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OAuth state"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in the request"})
		return
	}

	// Exchange code for a token.
	token, err := s.Config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Use the token to get user info from Google.
	client := s.Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}

	var googleUser GoogleUser
	if err := json.Unmarshal(body, &googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// Only accept users with a verified email.
	if !googleUser.VerifiedEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Google email not verified"})
		return
	}

	// Map the email to your user model.
	ctx := context.Background()
	user, err := s.UserStore.FindUserByEmail(ctx, googleUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if user == nil {
		// If no user exists, create a new user.
		user = &models.User{
			Username: googleUser.Name,
			// You may use the email as a unique identifier.
			// Adjust database fields as needed.
		}
		if err := s.UserStore.CreateUser(ctx, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// Generate a JWT token for the user.
	jwtToken, err := s.JWTService.GenerateToken(user.ID, googleUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Optionally pass the JWT token to the frontend by:
	// 1. Setting it as an HTTP-only cookie.
	// 2. Or redirecting back to the React app with the token in a URL fragment.
	// For this example, we'll set an HTTP-only cookie.
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Finally, redirect the user to the SPA.
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
