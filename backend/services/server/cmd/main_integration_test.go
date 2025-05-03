package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shivamkedia17/roshnii/services/server/internal/auth"
	"github.com/shivamkedia17/roshnii/services/server/internal/handlers"
	"github.com/shivamkedia17/roshnii/services/server/internal/middleware"
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
	"github.com/shivamkedia17/roshnii/shared/pkg/models"
	"github.com/shivamkedia17/roshnii/shared/pkg/storage"
)

var (
	testRouter *gin.Engine
	testDBPool *pgxpool.Pool
	testConfig *config.Config
)

// --- Test Setup (TestMain) ---

func TestMain(m *testing.M) {
	// --- Configuration ---
	// Ensure environment is set for testing (dev login)
	os.Setenv("ENVIRONMENT", "development")
	// Load config - Make sure it points to a TEST database!
	// You might need to adjust the path depending on where you run `go test`
	cfg, err := config.LoadConfig("../..") // Or "../.." if running from project root
	if err != nil {
		log.Fatalf("FATAL: Failed to load config for testing: %v", err)
	}
	// Override specific test settings if needed (e.g., ensure dev mode)
	cfg.Environment = "development"
	testConfig = cfg
	log.Printf("INFO: Running tests in environment: %s", cfg.Environment)
	log.Printf("INFO: Using Test Database URL (ensure this is a TEST DB): %s", cfg.PostgresURL) // Be careful!

	// --- Database Connection ---
	pool, err := pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to test database: %v", err)
	}
	// Ping DB
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		log.Fatalf("FATAL: Failed to ping test database: %v", err)
	}
	testDBPool = pool
	log.Println("INFO: Connected to test database successfully.")

	// --- Setup Application Dependencies (like in main.go) ---
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	database := &db.PostgresStore{Pool: testDBPool} // Use the concrete type with the test pool
	var storageService storage.BlobStorage
	if cfg.BlobStorageType == "local" {
		localStoragePath := cfg.LocalstoragePath
		if localStoragePath == "" {
			localStoragePath = "./uploads" // Default path
		}

		var err error
		storageService, err = storage.NewLocalStorage(localStoragePath)
		if err != nil {
			log.Fatalf("Failed to initialize local storage: %v", err)
		}
		log.Printf("Using local file storage at: %s", localStoragePath)
	} else {
		// For now, fall back to local storage if type is unrecognized
		log.Printf("Unrecognized storage type '%s', using local storage", cfg.BlobStorageType)
		storageService, _ = storage.NewLocalStorage("./uploads")
	}

	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTRefreshSecret, cfg.TokenDuration)
	googleOAuthService := auth.NewGoogleOAuthService(cfg, database, jwtService)

	imageHandler := handlers.NewImageHandler(database, storageService, cfg) // Pass the DB store interface

	// albumHandler := handlers.NewAlbumHandler(database, cfg)
	// tagHandler := handlers.NewTagHandler(database, cfg)
	// shareHandler := handlers.NewShareHandler(database, cfg)
	// searchHandler := handlers.NewSearchHandler(database, cfg)
	authMW := middleware.AuthMiddleware(jwtService)

	// --- Register Routes (like in main.go) ---
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
	api := router.Group("/api")
	{
		authRoutes := api.Group("/auth")
		{
			googleRoutes := authRoutes.Group("/google")
			{
				googleRoutes.GET("/login", googleOAuthService.HandleLogin)       // Will redirect in tests, maybe less useful
				googleRoutes.GET("/callback", googleOAuthService.HandleCallback) // Hard to test directly
				googleRoutes.POST("/logout", authMW, googleOAuthService.HandleLogout)
			}
			if cfg.Environment == "development" {
				devRoutes := authRoutes.Group("/dev")
				{
					// Directly use the handler func for dev login
					devRoutes.POST("/login", func(c *gin.Context) {
						var req struct {
							Email string `json:"email" binding:"required,email"`
							Name  string `json:"name"`
						}
						if err := c.ShouldBindJSON(&req); err != nil {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
							return
						}
						userName := req.Name
						if userName == "" {
							userName = "Dev User"
						}
						// Use the test DB instance
						user, err := database.FindOrCreateUserByEmail(c.Request.Context(), req.Email, userName, "dev")
						if err != nil {
							log.Printf("Test Dev Login Error: Failed to find/create user %s: %v", req.Email, err)
							c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login"})
							return
						}
						token, err := jwtService.GenerateToken(user)
						if err != nil {
							log.Printf("Test Dev Login Error: Failed to generate JWT for user %d: %v", user.ID, err)
							c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session"})
							return
						}
						c.JSON(http.StatusOK, gin.H{"token": token, "user_id": user.ID}) // Return user ID too for tests
					})
				}
			}
		}
		imageHandler.RegisterRoutes(api, authMW)
		// albumHandler.RegisterRoutes(api, authMW)  // Stubs, but register them
		// tagHandler.RegisterRoutes(api, authMW)    // Stubs
		// shareHandler.RegisterRoutes(api, authMW)  // Stubs
		// searchHandler.RegisterRoutes(api, authMW) // Stubs
		api.GET("/me", authMW, func(c *gin.Context) {
			claims := middleware.GetUserClaims(c)
			if claims == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"user_id": claims.UserID,
				"email":   claims.Email,
			})
		})
	}
	testRouter = router

	// --- Run Tests ---
	log.Println("INFO: Starting integration tests...")
	exitCode := m.Run()

	// --- Teardown ---
	log.Println("INFO: Cleaning up test resources...")
	testDBPool.Close()
	log.Println("INFO: Test database connection closed.")

	os.Exit(exitCode)
}

// --- Helper Functions ---

// clearTables truncates tables for test isolation. Call before tests needing a clean slate.
func clearTables(t *testing.T) {
	t.Helper()
	// Use CASCADE to handle foreign key constraints if necessary
	// RESTART IDENTITY resets sequence generators (like for user IDs)
	query := `TRUNCATE TABLE images, users RESTART IDENTITY CASCADE;`
	_, err := testDBPool.Exec(context.Background(), query)
	require.NoError(t, err, "Failed to truncate tables")
	// log.Println("DEBUG: Cleared tables: users, images")
}

// performRequest is a helper to make HTTP requests to the test router.
func performRequest(method, path string, body io.Reader, token string, contentType string) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	} else if body != nil {
		req.Header.Set("Content-Type", "application/json") // Default if body present
	}

	rr := httptest.NewRecorder()
	testRouter.ServeHTTP(rr, req)
	return rr, nil
}

// DevLoginResponse structure to parse the dev login response
type DevLoginResponse struct {
	Token  string        `json:"token"`
	UserID models.UserID `json:"user_id"`
	Error  string        `json:"error"` // Capture potential errors
}

// performDevLogin simulates the dev login and returns the token and user ID.
func performDevLogin(t *testing.T, email, name string) (string, models.UserID) {
	t.Helper()
	loginPayload := map[string]string{"email": email, "name": name}
	payloadBytes, _ := json.Marshal(loginPayload)

	rr, err := performRequest("POST", "/api/auth/dev/login", bytes.NewBuffer(payloadBytes), "", "application/json")
	require.NoError(t, err)

	// log.Printf("DEBUG: Dev Login Response Body for %s: %s", email, rr.Body.String())

	var respBody DevLoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	require.NoError(t, err, "Failed to unmarshal dev login response")

	// Use require for critical checks that prevent test continuation
	require.Equal(t, http.StatusOK, rr.Code, "Dev login failed with status %d. Body: %s", rr.Code, rr.Body.String())
	require.NotEmpty(t, respBody.Token, "Dev login response missing token")
	require.NotZero(t, respBody.UserID, "Dev login response missing user_id")
	require.Empty(t, respBody.Error, "Dev login response contained an error field: %s", respBody.Error) // Ensure no error field on success

	return respBody.Token, respBody.UserID
}

// ErrorResponse structure to parse standard error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// createDummyFile creates a temporary file for upload tests.
func createDummyFile(t *testing.T, filename, content string) string {
	t.Helper()
	tmpDir := t.TempDir() // Creates a temporary directory cleaned up after the test
	filePath := filepath.Join(tmpDir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err, "Failed to create dummy file")
	return filePath
}

// performUpload simulates uploading a file.
func performUpload(t *testing.T, token, fieldName, filePath, fileContentType string) (*httptest.ResponseRecorder, models.ImageMetadata) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Open the file
	file, err := os.Open(filePath)
	require.NoError(t, err, "Failed to open file for upload: %s", filePath)
	defer file.Close()

	// --- Create the form part with explicit headers ---
	// Create the MIME header for the part
	h := make(textproto.MIMEHeader)
	// Set Content-Disposition: form-data; name="file"; filename="imageA.jpg"
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			fieldName, filepath.Base(filePath)))
	// Set Content-Type: image/jpeg (or whatever is passed in)
	if fileContentType != "" { // Only set if provided
		h.Set("Content-Type", fileContentType)
	}

	// Create the part using the custom headers
	part, err := writer.CreatePart(h)
	require.NoError(t, err, "Failed to create form part with headers")
	// --- End Header Creation ---

	// Copy file content to the form part
	_, err = io.Copy(part, file)
	require.NoError(t, err, "Failed to copy file content to form part")

	// Close the multipart writer *before* making the request
	err = writer.Close()
	require.NoError(t, err, "Failed to close multipart writer")

	// Perform the request
	// IMPORTANT: Ensure the endpoint matches your ImageHandler registration (/api/upload)
	rr, err := performRequest("POST", "/api/upload", body, token, writer.FormDataContentType())
	require.NoError(t, err)

	// Check the response status code *before* trying to unmarshal
	require.Equal(t, http.StatusCreated, rr.Code, "Image upload failed with status %d. Body: %s", rr.Code, rr.Body.String())

	var uploadedMeta models.ImageMetadata
	err = json.Unmarshal(rr.Body.Bytes(), &uploadedMeta)
	require.NoError(t, err, "Failed to unmarshal upload response body: %s", rr.Body.String()) // Log body on error
	require.NotEmpty(t, uploadedMeta.ID, "Uploaded image metadata missing ID")
	// Add more checks if needed (e.g., returned content type)
	// require.Equal(t, fileContentType, uploadedMeta.ContentType, "Uploaded metadata content type mismatch")

	return rr, uploadedMeta
}

// --- Test Cases ---

func TestHealthCheck(t *testing.T) {
	rr, err := performRequest("GET", "/health", nil, "", "")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	var respBody map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	require.NoError(t, err)
	assert.Equal(t, "UP", respBody["status"])
}

func TestAuthenticationFlow(t *testing.T) {
	clearTables(t) // Ensure clean slate for user creation

	email := fmt.Sprintf("auth-tester-%d@example.com", time.Now().UnixNano())
	name := "Auth Tester"

	// 1. Dev Login Success
	t.Run("DevLoginSuccess", func(t *testing.T) {
		token, userID := performDevLogin(t, email, name)
		assert.NotEmpty(t, token)
		assert.NotZero(t, userID)
	})

	// 2. Dev Login Bad Request (Missing Email)
	t.Run("DevLoginBadRequest", func(t *testing.T) {
		loginPayload := map[string]string{"name": name} // Missing email
		payloadBytes, _ := json.Marshal(loginPayload)
		rr, err := performRequest("POST", "/api/auth/dev/login", bytes.NewBuffer(payloadBytes), "", "application/json")
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		// Check error message if desired
		var errResp ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "Invalid request body")
	})

	// 3. Get Me Unauthorized
	t.Run("GetMeUnauthorized", func(t *testing.T) {
		rr, err := performRequest("GET", "/api/me", nil, "", "") // No token
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	// 4. Get Me Authorized
	t.Run("GetMeAuthorized", func(t *testing.T) {
		token, userID := performDevLogin(t, email, name) // Login again to be sure

		rr, err := performRequest("GET", "/api/me", nil, token, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)

		var meResp map[string]interface{} // Use map for flexibility or define a struct
		err = json.Unmarshal(rr.Body.Bytes(), &meResp)
		require.NoError(t, err)
		// JSON numbers are decoded as float64 by default
		assert.Equal(t, float64(userID), meResp["user_id"])
		assert.Equal(t, email, meResp["email"])
	})

	// 5. Logout
	t.Run("Logout", func(t *testing.T) {
		token, _ := performDevLogin(t, email+"-logout", "Logout User") // Use different email

		rr, err := performRequest("POST", "/api/auth/google/logout", nil, token, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		// Check success message
		var logoutResp map[string]string
		json.Unmarshal(rr.Body.Bytes(), &logoutResp)
		assert.Equal(t, "Successfully logged out", logoutResp["message"])

		// Optional: Try accessing /me again with the same token - it *should* still work
		// because dev login doesn't rely on cookies being cleared by logout.
		// rr_post_logout, _ := performRequest("GET", "/api/me", nil, token, "")
		// assert.Equal(t, http.StatusOK, rr_post_logout.Code) // This demonstrates token is still valid
	})
}

func TestImageWorkflow_UserIsolation(t *testing.T) {
	clearTables(t) // *** Crucial for isolation ***

	// --- Setup Users ---
	email1 := fmt.Sprintf("user1-%d@example.com", time.Now().UnixNano())
	email2 := fmt.Sprintf("user2-%d@example.com", time.Now().UnixNano())
	token1, userID1 := performDevLogin(t, email1, "User One")
	token2, userID2 := performDevLogin(t, email2, "User Two")
	require.NotEqual(t, userID1, userID2, "User IDs should be different")

	// --- Initial State Check ---
	t.Run("InitialImageListEmpty", func(t *testing.T) {
		// User 1 checks images
		rr1, err1 := performRequest("GET", "/api/images", nil, token1, "")
		require.NoError(t, err1)
		require.Equal(t, http.StatusOK, rr1.Code)
		var images1 []models.ImageMetadata
		err1 = json.Unmarshal(rr1.Body.Bytes(), &images1)
		require.NoError(t, err1)
		assert.Empty(t, images1, "User 1 initial image list should be empty")

		// User 2 checks images
		rr2, err2 := performRequest("GET", "/api/images", nil, token2, "")
		require.NoError(t, err2)
		require.Equal(t, http.StatusOK, rr2.Code)
		var images2 []models.ImageMetadata
		err2 = json.Unmarshal(rr2.Body.Bytes(), &images2)
		require.NoError(t, err2)
		assert.Empty(t, images2, "User 2 initial image list should be empty")
	})

	// --- User 1 Uploads Images ---
	var imageAID, imageBID string // Store IDs for later checks
	t.Run("User1Uploads", func(t *testing.T) {
		// Use the fileContentType argument in performUpload
		dummyFileAPath := createDummyFile(t, "imageA.jpg", "jpeg content A")
		_, metaA := performUpload(t, token1, "file", dummyFileAPath, "image/jpeg") // <-- Pass content type
		assert.Equal(t, userID1, metaA.UserID)
		assert.Equal(t, "imageA.jpg", metaA.Filename)
		assert.Equal(t, "image/jpeg", metaA.ContentType) // Verify returned content type
		assert.NotEmpty(t, metaA.ID)
		imageAID = metaA.ID // Save ID

		dummyFileBPath := createDummyFile(t, "imageB.png", "png content B")
		_, metaB := performUpload(t, token1, "file", dummyFileBPath, "image/png") // <-- Pass content type
		assert.Equal(t, userID1, metaB.UserID)
		assert.Equal(t, "imageB.png", metaB.Filename)
		assert.Equal(t, "image/png", metaB.ContentType) // Verify returned content type
		assert.NotEmpty(t, metaB.ID)
		imageBID = metaB.ID // Save ID

		require.NotEqual(t, imageAID, imageBID, "Uploaded image IDs should be unique")
	})

	// --- Verify Image Lists After Upload (Isolation Check) ---
	t.Run("VerifyListsAfterUpload", func(t *testing.T) {
		// Add a small delay if tests are running extremely fast and DB operations might not be fully visible yet
		// time.Sleep(50 * time.Millisecond) // Usually not needed with transactions, but can help diagnose race conditions

		// User 1 lists images - should see A and B
		rr1, err1 := performRequest("GET", "/api/images", nil, token1, "")
		require.NoError(t, err1)
		require.Equal(t, http.StatusOK, rr1.Code)
		var images1 []models.ImageMetadata
		err1 = json.Unmarshal(rr1.Body.Bytes(), &images1)
		require.NoError(t, err1, "Failed to unmarshal User 1 image list: %s", rr1.Body.String()) // Log body on error
		require.Len(t, images1, 2, "User 1 should have 2 images. Body: %s", rr1.Body.String())   // Log body on error
		// Check if IDs are present (order might vary)
		foundA, foundB := false, false
		for _, img := range images1 {
			if img.ID == imageAID {
				foundA = true
				assert.Equal(t, "imageA.jpg", img.Filename)
				assert.Equal(t, "image/jpeg", img.ContentType)
			}
			if img.ID == imageBID {
				foundB = true
				assert.Equal(t, "imageB.png", img.Filename)
				assert.Equal(t, "image/png", img.ContentType)
			}
			assert.Equal(t, userID1, img.UserID, "Image listed by User 1 should belong to User 1")
		}
		assert.True(t, foundA, "Image A (ID: %s) not found in User 1 list", imageAID)
		assert.True(t, foundB, "Image B (ID: %s) not found in User 1 list", imageBID)

		// User 2 lists images - should see none
		rr2, err2 := performRequest("GET", "/api/images", nil, token2, "")
		require.NoError(t, err2)
		require.Equal(t, http.StatusOK, rr2.Code)
		var images2 []models.ImageMetadata
		err2 = json.Unmarshal(rr2.Body.Bytes(), &images2)
		require.NoError(t, err2)
		assert.Empty(t, images2, "User 2 image list should still be empty")
	})

	// --- Verify Get Specific Image (Isolation Check) ---
	t.Run("GetSpecificImageIsolation", func(t *testing.T) {
		require.NotEmpty(t, imageAID, "Cannot run GetSpecificImage test without imageAID") // This assertion should now pass

		// User 1 gets Image A - Success
		rr1, err1 := performRequest("GET", "/api/image/"+imageAID, nil, token1, "")
		require.NoError(t, err1)
		require.Equal(t, http.StatusOK, rr1.Code)
		var metaA models.ImageMetadata
		err1 = json.Unmarshal(rr1.Body.Bytes(), &metaA)
		require.NoError(t, err1)
		assert.Equal(t, imageAID, metaA.ID)
		assert.Equal(t, userID1, metaA.UserID)
		assert.Equal(t, "imageA.jpg", metaA.Filename)
		assert.Equal(t, "image/jpeg", metaA.ContentType) // Check content type here too

		// User 2 tries to get Image A - Failure (Not Found)
		rr2, err2 := performRequest("GET", "/api/image/"+imageAID, nil, token2, "")
		require.NoError(t, err2)
		// Ensure the GetImageByID store method correctly returns an error mapped to 404
		assert.Equal(t, http.StatusNotFound, rr2.Code, "User 2 should not be able to get User 1's image")
	})

	// --- Test Not Found Cases ---
	t.Run("GetNonExistentImage", func(t *testing.T) {
		// Use a valid UUID format but one that doesn't exist
		nonExistentUUID := "11111111-1111-1111-1111-111111111111"
		rr, err := performRequest("GET", "/api/image/"+nonExistentUUID, nil, token1, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rr.Code) // Assuming your handler returns 404 based on DB error
	})
}

func TestImageUploadValidations(t *testing.T) {
	clearTables(t)
	token, _ := performDevLogin(t, "validator@example.com", "Validator")

	t.Run("UploadMissingFileField", func(t *testing.T) {
		// Create multipart body *without* the file field
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		// Add some other field if needed by backend, otherwise just close
		writer.Close()

		rr, err := performRequest("POST", "/api/upload", body, token, writer.FormDataContentType())
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rr.Code) // Expecting 400 Bad Request
		var errResp ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "Missing or invalid 'file' field")
	})

	t.Run("UploadInvalidFileType", func(t *testing.T) {
		dummyTextPath := createDummyFile(t, "invalid.txt", "this is plain text")
		// Use performUpload helper, but expect failure
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		file, _ := os.Open(dummyTextPath)
		defer file.Close()
		part, _ := writer.CreateFormFile("file", filepath.Base(dummyTextPath))
		io.Copy(part, file)
		writer.Close()

		rr, err := performRequest("POST", "/api/upload", body, token, writer.FormDataContentType())
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var errResp ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "Unsupported file type") // Check based on your handler's error message
	})

	// Add test for file size limit if needed (requires creating a larger file)
}
func TestAlbumWorkflow(t *testing.T) {
	clearTables(t) // Ensure clean slate

	// Login first
	token, userID := performDevLogin(t, "album-test@example.com", "Album Tester")
	require.NotEmpty(t, token)
	require.NotZero(t, userID)

	// Variable to store albumID
	var albumID int64

	// Create a test album
	t.Run("CreateAlbum", func(t *testing.T) {
		albumPayload := map[string]string{
			"name":        "Integration Test Album",
			"description": "Created during integration testing",
		}
		payloadBytes, _ := json.Marshal(albumPayload)

		rr, err := performRequest("POST", "/api/albums", bytes.NewBuffer(payloadBytes), token, "application/json")
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, rr.Code, "Album creation failed with status %d. Body: %s", rr.Code, rr.Body.String())

		var album models.Album
		err = json.Unmarshal(rr.Body.Bytes(), &album)
		require.NoError(t, err)
		assert.Equal(t, "Integration Test Album", album.Name)
		assert.Equal(t, "Created during integration testing", album.Description)
		assert.Equal(t, userID, album.UserID)
		assert.NotZero(t, album.ID)

		// Store album ID for later tests
		albumID = album.ID
	})

	// List albums
	t.Run("ListAlbums", func(t *testing.T) {
		rr, err := performRequest("GET", "/api/albums", nil, token, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)

		var albums []models.Album
		err = json.Unmarshal(rr.Body.Bytes(), &albums)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(albums), 1, "Expected at least one album")

		// Find our test album
		found := false
		for _, album := range albums {
			if album.ID == albumID {
				found = true
				assert.Equal(t, "Integration Test Album", album.Name)
				break
			}
		}
		assert.True(t, found, "Created album not found in the list")
	})

	// Get specific album
	t.Run("GetAlbum", func(t *testing.T) {
		albumIDStr := fmt.Sprintf("%d", albumID)
		rr, err := performRequest("GET", "/api/albums/"+albumIDStr, nil, token, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)

		var album models.Album
		err = json.Unmarshal(rr.Body.Bytes(), &album)
		require.NoError(t, err)
		assert.Equal(t, albumID, album.ID)
		assert.Equal(t, "Integration Test Album", album.Name)
	})

	// Update album
	t.Run("UpdateAlbum", func(t *testing.T) {
		albumIDStr := fmt.Sprintf("%d", albumID)
		updatePayload := map[string]string{
			"name":        "Updated Test Album",
			"description": "This album was updated",
		}
		payloadBytes, _ := json.Marshal(updatePayload)

		rr, err := performRequest("PUT", "/api/albums/"+albumIDStr, bytes.NewBuffer(payloadBytes), token, "application/json")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)

		// Verify the update by fetching the album again
		rr, err = performRequest("GET", "/api/albums/"+albumIDStr, nil, token, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)

		var album models.Album
		err = json.Unmarshal(rr.Body.Bytes(), &album)
		require.NoError(t, err)
		assert.Equal(t, "Updated Test Album", album.Name)
		assert.Equal(t, "This album was updated", album.Description)
	})

	// Upload an image for album/image tests
	var imageID string
	t.Run("UploadImageForAlbum", func(t *testing.T) {
		dummyImagePath := createDummyFile(t, "album_test_image.jpg", "dummy image content")
		_, meta := performUpload(t, token, "file", dummyImagePath, "image/jpeg")
		imageID = meta.ID
		require.NotEmpty(t, imageID)
	})

	// Add image to album
	t.Run("AddImageToAlbum", func(t *testing.T) {
		albumIDStr := fmt.Sprintf("%d", albumID)
		addImagePayload := map[string]string{
			"image_id": imageID,
		}
		payloadBytes, _ := json.Marshal(addImagePayload)

		rr, err := performRequest("POST", "/api/albums/"+albumIDStr+"/images", bytes.NewBuffer(payloadBytes), token, "application/json")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)
	})

	// List images in album
	t.Run("ListImagesInAlbum", func(t *testing.T) {
		albumIDStr := fmt.Sprintf("%d", albumID)
		rr, err := performRequest("GET", "/api/albums/"+albumIDStr+"/images", nil, token, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)

		var images []models.ImageMetadata
		err = json.Unmarshal(rr.Body.Bytes(), &images)
		require.NoError(t, err)
		assert.Equal(t, 1, len(images), "Expected exactly one image in album")
		assert.Equal(t, imageID, images[0].ID)
	})

	// Remove image from album
	t.Run("RemoveImageFromAlbum", func(t *testing.T) {
		albumIDStr := fmt.Sprintf("%d", albumID)
		rr, err := performRequest("DELETE", "/api/albums/"+albumIDStr+"/images/"+imageID, nil, token, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)

		// Verify image is removed
		rr, err = performRequest("GET", "/api/albums/"+albumIDStr+"/images", nil, token, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)

		var images []models.ImageMetadata
		err = json.Unmarshal(rr.Body.Bytes(), &images)
		require.NoError(t, err)
		assert.Equal(t, 0, len(images), "Expected no images in album after removal")
	})

	// Delete album
	t.Run("DeleteAlbum", func(t *testing.T) {
		albumIDStr := fmt.Sprintf("%d", albumID)
		rr, err := performRequest("DELETE", "/api/albums/"+albumIDStr, nil, token, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)

		// Verify album is deleted
		rr, err = performRequest("GET", "/api/albums/"+albumIDStr, nil, token, "")
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, rr.Code, "Album should be deleted")
	})
}
