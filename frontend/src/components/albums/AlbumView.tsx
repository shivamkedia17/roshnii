import { useState, useEffect } from "react";
import {
  useAlbum,
  useAlbumImages,
  useUpdateAlbum,
  useRemoveImageFromAlbum,
} from "@/hooks/useAlbumQueries";
import { PhotoItem } from "../photos/PhotoItem";
import { PhotoModal } from "../photos/PhotoModal";
import "@/css/AlbumView.css";

type AlbumViewProps = {
  albumId: string;
  onBack: () => void;
};

export function AlbumView({ albumId, onBack }: AlbumViewProps) {
  const id = parseInt(albumId, 10);
  const [selectedPhoto, setSelectedPhoto] = useState<string | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [editName, setEditName] = useState("");
  const [editDescription, setEditDescription] = useState("");
  const [confirmRemoveImage, setConfirmRemoveImage] = useState<string | null>(
    null,
  );

  // Fetch album details and images
  const {
    data: album,
    isLoading: albumLoading,
    error: albumError,
  } = useAlbum(id);

  const {
    data: images,
    isLoading: imagesLoading,
    error: imagesError,
  } = useAlbumImages(id);

  // Update album mutation
  const updateAlbumMutation = useUpdateAlbum();

  // Remove image from album mutation
  const removeImageMutation = useRemoveImageFromAlbum();

  // Set initial form values when album data loads
  useEffect(() => {
    if (album) {
      setEditName(album.name);
      setEditDescription(album.description || "");
    }
  }, [album]);

  // Handle form submission for album update
  const handleUpdateAlbum = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editName.trim()) return;

    try {
      await updateAlbumMutation.mutateAsync({
        id,
        data: {
          name: editName.trim(),
          description: editDescription.trim() || undefined,
        },
      });
      setIsEditing(false);
    } catch (err) {
      console.error("Failed to update album:", err);
    }
  };

  // Handle removing an image from the album
  const handleRemoveImage = async (imageId: string) => {
    try {
      await removeImageMutation.mutateAsync({ albumId: id, imageId });
      setConfirmRemoveImage(null);
    } catch (err) {
      console.error("Failed to remove image from album:", err);
    }
  };

  // Loading state
  if (albumLoading || imagesLoading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner"></div>
        <p>Loading album...</p>
      </div>
    );
  }

  // Error state
  if (albumError || imagesError) {
    return (
      <div className="error-container">
        <p>
          Error loading album:
          {albumError instanceof Error
            ? albumError.message
            : imagesError instanceof Error
              ? imagesError.message
              : "Unknown error"}
        </p>
        <button onClick={onBack} className="back-button">
          Back to Albums
        </button>
      </div>
    );
  }

  // Album not found
  if (!album) {
    return (
      <div className="error-container">
        <p>Album not found</p>
        <button onClick={onBack} className="back-button">
          Back to Albums
        </button>
      </div>
    );
  }

  return (
    <div className="album-view">
      <div className="album-view-header">
        <button onClick={onBack} className="back-button">
          ← Back to Albums
        </button>

        {isEditing ? (
          <form onSubmit={handleUpdateAlbum} className="edit-album-form">
            <div className="form-group">
              <label htmlFor="edit-name">Album Name</label>
              <input
                id="edit-name"
                type="text"
                value={editName}
                onChange={(e) => setEditName(e.target.value)}
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="edit-description">Description</label>
              <textarea
                id="edit-description"
                value={editDescription}
                onChange={(e) => setEditDescription(e.target.value)}
                rows={2}
              />
            </div>
            <div className="form-buttons">
              <button
                type="button"
                className="cancel-button"
                onClick={() => {
                  setIsEditing(false);
                  // Reset form to original values
                  setEditName(album.name);
                  setEditDescription(album.description || "");
                }}
              >
                Cancel
              </button>
              <button
                type="submit"
                className="save-button"
                disabled={updateAlbumMutation.isPending || !editName.trim()}
              >
                {updateAlbumMutation.isPending ? "Saving..." : "Save Changes"}
              </button>
            </div>
            {updateAlbumMutation.error && (
              <div className="error">
                {updateAlbumMutation.error instanceof Error
                  ? updateAlbumMutation.error.message
                  : "Failed to update album"}
              </div>
            )}
          </form>
        ) : (
          <div className="album-info">
            <div className="album-title-section">
              <h1>{album.name}</h1>
              <button
                onClick={() => setIsEditing(true)}
                className="edit-button"
              >
                Edit
              </button>
            </div>
            {album.description && (
              <p className="album-description">{album.description}</p>
            )}
            <p className="album-date">
              Created {new Date(album.created_at).toLocaleDateString()}
              {album.created_at !== album.updated_at &&
                ` · Updated ${new Date(album.updated_at).toLocaleDateString()}`}
            </p>
          </div>
        )}
      </div>

      <div className="album-content">
        {!images || images.length === 0 ? (
          <div className="no-photos">
            <div className="no-photos-icon">🖼️</div>
            <p>This album is empty</p>
            <button
              onClick={() => (window.location.hash = "#upload")}
              className="upload-now-button"
            >
              Upload Photos
            </button>
          </div>
        ) : (
          <>
            <div className="album-photos-header">
              <h2>Photos</h2>
              <p className="photo-count">{images.length} photos</p>
            </div>
            <div className="gallery-grid">
              {images.map((photo) => (
                <div key={photo.id} className="album-photo-container">
                  <PhotoItem
                    photo={photo}
                    onClick={() => setSelectedPhoto(photo.id)}
                  />
                  <button
                    className="remove-from-album-button"
                    onClick={() => setConfirmRemoveImage(photo.id)}
                    title="Remove from album"
                  >
                    ×
                  </button>

                  {/* Confirmation dialog for removing image */}
                  {confirmRemoveImage === photo.id && (
                    <div className="remove-confirm-dialog">
                      <p>Remove from album?</p>
                      {removeImageMutation.error && (
                        <div className="error">
                          {removeImageMutation.error instanceof Error
                            ? removeImageMutation.error.message
                            : "Failed to remove image"}
                        </div>
                      )}
                      <div className="confirm-buttons">
                        <button
                          onClick={() => setConfirmRemoveImage(null)}
                          className="cancel-button"
                          disabled={removeImageMutation.isPending}
                        >
                          Cancel
                        </button>
                        <button
                          onClick={() => handleRemoveImage(photo.id)}
                          className="confirm-remove-button"
                          disabled={removeImageMutation.isPending}
                        >
                          {removeImageMutation.isPending
                            ? "Removing..."
                            : "Remove"}
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          </>
        )}
      </div>

      {selectedPhoto && (
        <PhotoModal
          photoId={selectedPhoto}
          onClose={() => setSelectedPhoto(null)}
          albumContext={{
            albumId: id,
            onRemoveFromAlbum: (photoId) => {
              handleRemoveImage(photoId);
              setSelectedPhoto(null);
            },
          }}
        />
      )}
    </div>
  );
}
