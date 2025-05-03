import { useState } from "react";
import { usePhoto, useDeletePhoto } from "@/hooks/usePhotoQueries";
import "@/css/PhotoModal.css"; // Create this file next
import { AddToAlbumModal } from "../albums/AddToAlbumModal";

// Update the props type
export type PhotoModalProps = {
  photoId: string;
  onClose: () => void;
  albumContext?: {
    albumId: number;
    onRemoveFromAlbum: (photoId: string) => void;
  };
};

export function PhotoModal({
  photoId,
  onClose,
  albumContext,
}: PhotoModalProps) {
  const [isDeleteConfirmOpen, setIsDeleteConfirmOpen] = useState(false);
  const [showAddToAlbumModal, setShowAddToAlbumModal] = useState(false);
  const { data: photo, isLoading, error } = usePhoto(photoId);
  const deleteMutation = useDeletePhoto();

  if (isLoading) return <div className="photo-modal loading">Loading...</div>;
  if (error)
    return (
      <div className="photo-modal error">
        Error loading photo:{" "}
        {error instanceof Error ? error.message : "Unknown error"}
      </div>
    );
  if (!photo) return null;

  const handleDelete = async () => {
    try {
      await deleteMutation.mutateAsync(photoId);
      onClose(); // Close modal after successful deletion
    } catch (err) {
      console.error("Failed to delete photo:", err);
      // Error will be shown via deleteMutation.error
    }
  };

  return (
    <div className="photo-modal-overlay" onClick={onClose}>
      <div className="photo-modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>{photo.filename || "Photo"}</h3>
          <button className="close-button" onClick={onClose}>
            ×
          </button>
        </div>

        <div className="modal-body">
          <div className="photo-container">
            <img
              src={`/api/image/${photo.id}/download`}
              alt={photo.filename || "Photo"}
              className="full-size-photo"
            />
          </div>

          <div className="photo-details">
            <div className="detail-row">
              <span className="detail-label">Uploaded:</span>
              <span className="detail-value">
                {new Date(photo.created_at).toLocaleString()}
              </span>
            </div>
            <div className="detail-row">
              <span className="detail-label">Type:</span>
              <span className="detail-value">{photo.content_type}</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">Size:</span>
              <span className="detail-value">{formatFileSize(photo.size)}</span>
            </div>

            {photo.width && photo.height && (
              <div className="detail-row">
                <span className="detail-label">Dimensions:</span>
                <span className="detail-value">
                  {photo.width} × {photo.height}
                </span>
              </div>
            )}
          </div>
        </div>

        <div className="modal-footer">
          {/* Add a new button for adding to album - if not in album context */}
          {!albumContext && (
            <button
              className="add-to-album-button"
              onClick={() => setShowAddToAlbumModal(true)}
              disabled={deleteMutation.isPending}
            >
              Add to Album
            </button>
          )}

          {/* If in album context, show remove from album button */}
          {albumContext && (
            <button
              className="remove-from-album-button"
              onClick={() => albumContext.onRemoveFromAlbum(photoId)}
              disabled={deleteMutation.isPending}
            >
              Remove from Album
            </button>
          )}

          <button
            className="delete-button"
            onClick={() => setIsDeleteConfirmOpen(true)}
            disabled={deleteMutation.isPending}
          >
            Delete Photo
          </button>
        </div>

        {/* Delete confirmation dialog */}
        {isDeleteConfirmOpen && (
          <div className="delete-confirm-dialog">
            <h4>Delete Photo?</h4>
            <p>This action cannot be undone.</p>

            {deleteMutation.error && (
              <div className="error">
                {deleteMutation.error instanceof Error
                  ? deleteMutation.error.message
                  : "Failed to delete photo"}
              </div>
            )}

            <div className="confirm-buttons">
              <button
                onClick={() => setIsDeleteConfirmOpen(false)}
                className="cancel-button"
                disabled={deleteMutation.isPending}
              >
                Cancel
              </button>
              <button
                onClick={handleDelete}
                className="confirm-delete-button"
                disabled={deleteMutation.isPending}
              >
                {deleteMutation.isPending ? "Deleting..." : "Confirm Delete"}
              </button>
            </div>
          </div>
        )}

        {showAddToAlbumModal && (
          <AddToAlbumModal
            imageId={photoId}
            onClose={() => setShowAddToAlbumModal(false)}
          />
        )}
      </div>
    </div>
  );
}

// Helper function to format file size
function formatFileSize(bytes: number): string {
  if (bytes < 1024) return bytes + " B";
  else if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
  else return (bytes / (1024 * 1024)).toFixed(1) + " MB";
}
