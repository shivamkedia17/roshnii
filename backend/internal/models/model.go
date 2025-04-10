package models

import "time"

// User represents a registered user in the system.
type User struct {
	// ID is the unique identifier for the user (typically auto-generated by the database).
	ID int64 `json:"id"`

	// Username is the unique name used for login.
	Username string `json:"username"`

	// PasswordHash is the hashed password for storage.
	// The "-" tag prevents this field from being included in JSON responses.
	PasswordHash string `json:"-"`

	// CreatedAt timestamp indicates when the user account was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt timestamp indicates when the user account was last updated.
	// Using a pointer allows us to distinguish between a zero time and a time not being set.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// ImageMetadata holds information about an uploaded image.
type ImageMetadata struct {
	// ID is the unique identifier for the image (e.g., a UUID generated by the application).
	ID string `json:"id"`

	// UserID links the image to the user who uploaded it.
	UserID int64 `json:"user_id"`

	// Filename is the original name of the uploaded file.
	Filename string `json:"filename"`

	// StoragePath is the internal path or key used to retrieve the image from blob storage.
	// This is typically not exposed via the API.
	StoragePath string `json:"-"`

	// ContentType is the MIME type of the image (e.g., "image/jpeg", "image/png").
	ContentType string `json:"content_type"`

	// Size is the file size in bytes.
	Size int64 `json:"size"`

	// Width is the width of the image in pixels. (May be populated later if not available at upload)
	Width int `json:"width,omitempty"`

	// Height is the height of the image in pixels. (May be populated later if not available at upload)
	Height int `json:"height,omitempty"`

	// PHash is the perceptual hash of the image, used for deduplication.
	// Using a pointer allows us to represent the state where the hash hasn't been calculated yet.
	// Typically not exposed directly via standard APIs unless specifically needed.
	PHash *string `json:"-"`

	// IsDuplicate indicates if this image is considered a duplicate of another.
	IsDuplicate bool `json:"is_duplicate"`

	// OriginalID points to the ID of the original image if this one is a duplicate.
	// Using a pointer allows us to represent non-duplicates (nil value).
	OriginalID *string `json:"original_id,omitempty"`

	// UploadedAt timestamp indicates when the image was successfully uploaded and metadata created.
	UploadedAt time.Time `json:"uploaded_at"`

	// TakenAt represents the date and time the photo was actually taken, often extracted from EXIF data.
	// Using a pointer as it might not always be available.
	TakenAt *time.Time `json:"taken_at,omitempty"`

	// Location represents geo-coordinates extracted from EXIF data. Could be a nested struct or separate fields.
	// Example:
	// Latitude *float64 `json:"latitude,omitempty"`
	// Longitude *float64 `json:"longitude,omitempty"`

	// NeedsEmbedding indicates if the embedding needs to be generated (or regenerated).
	// Useful for background job processing. Could also be managed via message queues.
	// NeedsEmbedding bool `json:"-"`

	// NeedsClustering indicates if the image needs to be considered in the next face clustering run.
	// NeedsClustering bool `json:"-"`
}

// Potential other structs could go here later, e.g.:
// - FaceCluster
// - Album
// - APIRequest/Response specific structs if needed
