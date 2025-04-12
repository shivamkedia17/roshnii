
Below is a complete set of steps and sample code to implement Google OAuth login in your Roshnii app. In this example, the backend (built with Go and Gin) will handle the OAuth flow; once the user logs in with Google, the backend extracts the user's email, maps it to your database’s user_id, and issues a JWT token for the React SPA to use.

──────────────────────────────
Step 1. Set Up Google OAuth Credentials
──────────────────────────────
1. Go to the Google Cloud Console (https://console.cloud.google.com/).
2. Create or select a project.
3. Go to APIs & Services > Credentials.
4. Create OAuth client credentials for a “Web application.”
5. Add authorized redirect URIs (for example, for development:
   http://localhost:8080/api/auth/google/callback).
6. Save the Client ID and Client Secret.
7. Add these values to your backend’s configuration (for example, in environment variables or your .env file).

──────────────────────────────
Step 2. Create the OAuth Service (Backend)
──────────────────────────────
Create a file such as internal/auth/google_oauth.go. This file will include:
• An OAuth2 configuration using your Google credentials.
• A method to generate and verify the OAuth state parameter.
• Handlers for both initiating the OAuth login and handling the callback.

Below is an example implementation:

--------------------------------------------------
File: internal/auth/google_oauth.go
--------------------------------------------------
package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shivamkedia17/roshnii-backend/internal/config"
	"github.com/shivamkedia17/roshnii-backend/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

──────────────────────────────
Step 3. Wire Up the OAuth Endpoints in the Router
──────────────────────────────
In your backend router (for example, using Gin), add routes that call the methods above.

For example, in your main router file (e.g., internal/server/router.go):

--------------------------------------------------
router.GET("/api/auth/google/login", googleOAuthService.OAuthLogin)
router.GET("/api/auth/google/callback", googleOAuthService.OAuthCallback)
--------------------------------------------------

Make sure your configuration (including cfg.GoogleClientID and cfg.GoogleClientSecret) is loaded before initializing the GoogleOAuthService.

──────────────────────────────
Step 4. Update Your User Database Model
──────────────────────────────
Ensure your user table (or equivalent) has a unique field for the email address. The mapping is as follows:
• OAuth callback returns GoogleUser.Email.
• In your application’s user table, this email becomes the unique identifier (or is used to lookup associated user_id).

In models/model.go (or your user model file), you might have:
--------------------------------------------------
type User struct {
    ID       int64  `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`  // Add this field if not already present.
    // ... other fields ...
}
--------------------------------------------------

──────────────────────────────
Step 5. Frontend: Triggering the Login Flow
──────────────────────────────
On the frontend, update your Login component (components/auth/Login.tsx) so that the “Sign in with Google” button redirects the user to the backend API endpoint that initiates OAuth.

For example:
--------------------------------------------------
export function Login() {
  const handleGoogleLogin = () => {
    // This URL will trigger your backend's "/api/auth/google/login" endpoint.
    window.location.href = "/api/auth/google/login";
  };

  return (
    <div className="login-container">
      <div className="login-card">
        <h1>Roshnii</h1>
        <p>Store and organize your memories</p>
        <button className="google-login-btn" onClick={handleGoogleLogin}>
          Sign in with Google
        </button>
      </div>
    </div>
  );
}
--------------------------------------------------
Step 6. Finalize Session Handling
──────────────────────────────
1. Upon receiving the callback and setting the auth_token cookie, the SPA should check for the token existence (either via a dedicated API endpoint or from cookies) and update the authentication context accordingly.
2. On each authenticated API call, send the JWT (the cookie is automatically included if using cookie-based auth) so that your backend can validate the session.

──────────────────────────────
Summary
──────────────────────────────
• Configure Google OAuth credentials and add them to your backend configuration.
• Implement an OAuth service in Go that handles login and callback, mapping the Google email to your user record.
• Wire the endpoints in your Gin router.
• On the frontend, have the login button redirect to the backend OAuth endpoint.
• After OAuth, use a JWT token (typically stored in an HTTP-only cookie or passed via the URL fragment) for maintaining the session in your React SPA.

By following these steps, your application will leverage Google OAuth for login while ensuring that each user’s email is mapped correctly to their database record, and the user receives a JWT for subsequent API requests.
Which OAuth grant type have you given instructions for? Is it a good idea to use in a single page app wiht an API-based backend?
The instructions provided implement the standard Authorization Code Grant Flow. In this flow, your backend acts as a confidential client that performs the code exchange securely with Google (using your client secret), then generates its own JWT for the SPA.

For a single-page app with an API-based backend, this approach is good—provided that:

• The sensitive parts of the OAuth flow (code exchange, handling the client secret, etc.) occur on the backend rather than in the SPA.

• You secure the communication and token storage (for example, using HTTP-only cookies for the JWT).

If you were to build a public client (a pure SPA) without a secure backend handling secrets, you’d need to use the Authorization Code Flow with PKCE (Proof Key for Code Exchange) instead. PKCE protects against interception attacks when the client cannot safely store a secret.

In summary, using the standard Authorization Code Grant Flow as shown is a good idea when your backend is an API server that securely handles the OAuth code exchange. It keeps your client secret secure and makes it easier to manage user sessions.
