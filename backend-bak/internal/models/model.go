package models

type UserID = int64
type ImageID = string

// User represents a registered user in the system.
type User struct {
	ID           UserID `json:"id"` // Auto-generated unique user id
	Email        string `json:"email"`
	Name         string `json:"username"`      // User's display name
	AuthProvider string `json:"auth_provider"` // OAuth Identity Provider
	// The "-" tag prevents this field from being included in JSON responses.
	PasswordHash string `json:"-"`
}

// ImageMetadata holds information about an uploaded image.
type ImageMetadataCore struct {
	ID          ImageID `json:"id"` // Auto-generated unique image id
	UserID      UserID  `json:"user_id"`
	Filename    string  `json:"filename"`     // o.g name of the uploaded file
	StoragePath string  `json:"-"`            // blob-storage path
	ContentType string  `json:"content_type"` // MIME type of the image (e.g., "image/jpeg", "image/png").
	Size        int64   `json:"size"`
	Width       int     `json:"width,omitempty"`  // (May be populated later if not available at upload)
	Height      int     `json:"height,omitempty"` //(May be populated later if not available at upload)
}

//

// type ImageMetadataAdditional struct {
// 	// PHash is the perceptual hash of the image, used for deduplication.
// 	// Using a pointer allows us to represent the state where the hash hasn't been calculated yet.
// 	PHash       *string `json:"-"`
// 	IsDuplicate bool    `json:"is_duplicate"`

// 	// OriginalID points to the ID of the original image if this one is a duplicate.
// 	// Using a pointer allows us to represent non-duplicates (nil value).
// 	OriginalID *string `json:"original_id,omitempty"`

// 	// UploadedAt timestamp indicates when the image was successfully uploaded and metadata created.
// 	UploadedAt time.Time `json:"uploaded_at"`

// 	// TakenAt represents the date and time the photo was actually taken, often extracted from EXIF data.
// 	// Using a pointer as it might not always be available.
// 	TakenAt *time.Time `json:"taken_at,omitempty"`

// 	// Location represents geo-coordinates extracted from EXIF data. Could be a nested struct or separate fields.
// 	// Example:
// 	// Latitude *float64 `json:"latitude,omitempty"`
// 	// Longitude *float64 `json:"longitude,omitempty"`

// 	// NeedsEmbedding indicates if the embedding needs to be generated (or regenerated).
// 	// Useful for background job processing. Could also be managed via message queues.
// 	// NeedsEmbedding bool `json:"-"`

// 	// NeedsClustering indicates if the image needs to be considered in the next face clustering run.
// 	// NeedsClustering bool `json:"-"`
// }

// Potential other structs could go here later, e.g.:
// - FaceCluster
// - Album
// - APIRequest/Response specific structs if needed
