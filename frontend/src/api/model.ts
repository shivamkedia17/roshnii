// Type aliases
export type UserID = string; // UUID
export type ImageID = string; // UUID
export type AlbumID = string; // UUID

// User represents a registered user in the system.
export type User = {
  id: UserID;
  email: string;
  name: string;
  picture_url?: string;
};

// ImageMetadata holds information about an uploaded image.
export type ImageMetadata = {
  id: ImageID;
  user_id: UserID;
  filename: string; // Original filename
  content_type: string; // MIME type
  size: number; // Size in bytes
  width?: number;
  height?: number;
  created_at: Date;
  updated_at: Date;
};

// Album represents a collection of images grouped by a user.
export type Album = {
  id: AlbumID;
  user_id: UserID;
  name: string;
  description: string;
  created_at: Date;
  updated_at: Date;
};

// AlbumImage links images to albums (many-to-many).
export type AlbumImage = {
  album_id: AlbumID;
  image_id: ImageID;
};

export type ServerMessage = {
  message: string;
};
