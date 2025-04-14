package models

import "time" // Added	// Added	// Added	// Added	// Added	// Added	// Added	// Added

type (
	UserID  = int64
	ImageID = string // Consider using UUID later if needed
)

// User represents a registered user in the system.
type User struct {
	ID           UserID    `json:"id" db:"id"`       // Auto-generated unique user id
	GoogleID     string    `json:"-" db:"google_id"` // Store Google's unique ID
	Email        string    `json:"email" db:"email"`
	Name         string    `json:"name" db:"name"`               // User's display name from Google
	PictureURL   string    `json:"picture_url" db:"picture_url"` // Profile picture URL from Google
	AuthProvider string    `json:"-" db:"auth_provider"`         // e.g., "google"
	CreatedAt    time.Time `json:"-" db:"created_at"`
	UpdatedAt    time.Time `json:"-" db:"updated_at"`
	// PasswordHash string    `json:"-" db:"password_hash"` // Not needed for pure OAuth
}

// ImageMetadata holds information about an uploaded image.
type ImageMetadata struct {
	ID          ImageID   `json:"id" db:"id"` // UUID or other unique ID
	UserID      UserID    `json:"user_id" db:"user_id"`
	Filename    string    `json:"filename" db:"filename"`         // Original filename
	StoragePath string    `json:"-" db:"storage_path"`            // Path in blob storage
	ContentType string    `json:"content_type" db:"content_type"` // MIME type
	Size        int64     `json:"size" db:"size"`                 // Size in bytes
	Width       int       `json:"width,omitempty" db:"width"`
	Height      int       `json:"height,omitempty" db:"height"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Add other fields later as needed (phash, taken_at, etc.)
}

type GoogleUser struct {
	ID            string `json:"sub"`            // Google's unique subject identifier
	Email         string `json:"email"`          // User's email address
	VerifiedEmail bool   `json:"email_verified"` // Whether Google has verified the email
	Name          string `json:"name"`           // User's full name
	GivenName     string `json:"given_name"`     // First name
	FamilyName    string `json:"family_name"`    // Last name
	Picture       string `json:"picture"`        // URL to profile picture
	Locale        string `json:"locale"`         // e.g., "en"
}
