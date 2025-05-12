import { useState, useEffect } from "react";
import { useListImage, useGetImage, useDeleteImage } from "@/hooks/useImages";
import { useRemoveImageFromAlbum } from "@/hooks/useAlbums";
import { AddToAlbumModal } from "../albums/AddToAlbumModal";
import "@/css/PhotoModal.css";

export type PhotoModalProps = {
  imageId: string;
  albumId?: string;
  onClose: () => void;
};

export function PhotoModal({ imageId, albumId, onClose }: PhotoModalProps) {
  const [isDeleteConfirmOpen, setIsDeleteConfirmOpen] = useState(false);
  const [showAddToAlbumModal, setShowAddToAlbumModal] = useState(false);

  const {
    data: imageMetadata,
    isLoading: isMetadataLoading,
    isFetching: isMetadataFetching,
    error: metadataError,
  } = useListImage(imageId);

  const {
    data: imageBlob,
    isLoading: isBlobLoading,
    isFetching: isBlobFetching,
    error: blobError,
  } = useGetImage(imageId);

  // Mutation Functions
  const deleteImage = useDeleteImage();
  const removeFromAlbum = useRemoveImageFromAlbum();

  const handleDelete = () => {
    deleteImage(imageId);
    onClose(); // Close modal after triggering deletion
  };

  const handleRemoveFromAlbum = () => {
    if (!albumId) return;

    removeFromAlbum({ albumId, imageId });
    onClose();
  };

  // Derived states
  const isLoading = isMetadataLoading || isBlobLoading;
  const isFetching = isMetadataFetching || isBlobFetching;
  const error = metadataError || blobError;

  // Create object URL for the blob and ensure proper cleanup
  const [imageBlobURL, setImageBlobURL] = useState<string>("");
  const [imageLoaded, setImageLoaded] = useState(false);
  const [imageDimensions, setImageDimensions] = useState({
    width: 0,
    height: 0,
  });

  // Set up and clean up object URL when blob changes
  useEffect(() => {
    if (imageBlob) {
      const url = URL.createObjectURL(imageBlob);
      setImageBlobURL(url);
      setImageLoaded(false);

      // revoke the object URL when component unmounts or blob changes
      return () => {
        URL.revokeObjectURL(url);
      };
    }
  }, [imageBlob]);

  // Loading State
  if (isLoading || isFetching) {
    return (
      <div className="photo-modal-overlay">
        <div className="photo-modal loading">
          <div className="loading-spinner"></div>
          <p>Loading photo...</p>
        </div>
      </div>
    );
  }

  // Error State
  if (error || !imageMetadata) {
    return (
      <div className="photo-modal-overlay">
        <div className="photo-modal error">
          <h3>Error</h3>
          <p>
            {error instanceof Error
              ? error.message
              : "Failed to load photo information"}
          </p>
          <button onClick={onClose}>Close</button>
        </div>
      </div>
    );
  }

  // All is well
  return (
    <div className="photo-modal-overlay" onClick={onClose}>
      <div className="photo-modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>{imageMetadata.filename || "Photo"}</h3>
          <button className="close-button" onClick={onClose}>
            ×
          </button>
        </div>

        <div className="modal-body">
          <div className="photo-container">
            {imageBlobURL && (
              <img
                src={imageBlobURL}
                alt={imageMetadata.filename || "Photo"}
                className="full-size-photo"
                onLoad={(e) => {
                  const img = e.target as HTMLImageElement;
                  setImageDimensions({
                    width: img.naturalWidth,
                    height: img.naturalHeight,
                  });
                  setImageLoaded(true);
                }}
              />
            )}
          </div>

          <div className="photo-details">
            <div className="detail-row">
              <span className="detail-label">Uploaded:</span>
              <span className="detail-value">
                {new Date(imageMetadata.created_at).toLocaleString()}
              </span>
            </div>
            <div className="detail-row">
              <span className="detail-label">Type:</span>
              <span className="detail-value">{imageMetadata.content_type}</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">Size:</span>
              <span className="detail-value">
                {formatFileSize(imageMetadata.size)}
              </span>
            </div>

            <div className="detail-row">
              <span className="detail-label">Dimensions:</span>
              <span className="detail-value">
                {imageLoaded && imageDimensions.width > 0
                  ? `${imageDimensions.width} × ${imageDimensions.height}`
                  : imageMetadata.width && imageMetadata.height
                    ? `${imageMetadata.width} × ${imageMetadata.height}`
                    : "Unavailable"}
              </span>
            </div>
          </div>
        </div>

        <div className="modal-footer">
          {/* Only show Add to Album button when not in album context */}
          {!albumId ? (
            <button
              className="add-to-album-button"
              onClick={() => setShowAddToAlbumModal(true)}
              // disabled={deleteImage}
            >
              Add to Album
            </button>
          ) : (
            <button
              className="remove-from-album-button"
              onClick={handleRemoveFromAlbum}
              // disabled={removeFromAlbum.isPending}
            >
              {/* {removeFromAlbum.isPending ? "Removing..." : "Remove from Album"} */}
              {"Remove from Album"}
            </button>
          )}

          <button
            className="delete-button"
            onClick={() => setIsDeleteConfirmOpen(true)}
            // disabled={deleteImage.isPending}
          >
            Delete Photo
          </button>
        </div>

        {/* Delete confirmation dialog */}
        {isDeleteConfirmOpen && (
          <div className="delete-confirm-dialog">
            <h4>Delete Photo?</h4>
            <p>This action cannot be undone.</p>

            {/* {deleteImage.error && (
              <div className="error">
                {deleteImage.error instanceof Error
                  ? deleteImage.error.message
                  : "Failed to delete photo"}
              </div>
            )} */}

            <div className="confirm-buttons">
              <button
                onClick={() => setIsDeleteConfirmOpen(false)}
                className="cancel-button"
                // disabled={deleteImage.isPending}
              >
                Cancel
              </button>
              <button
                onClick={handleDelete}
                className="confirm-delete-button"
                // disabled={deleteImage.isPending}
              >
                {/* {deleteImage.isPending ? "Deleting..." : "Confirm Delete"} */}
                {"Confirm Delete"}
              </button>
            </div>
          </div>
        )}

        {/* Add to Album Modal */}
        {showAddToAlbumModal && (
          <AddToAlbumModal
            imageId={imageId}
            onClose={() => setShowAddToAlbumModal(false)}
            onSuccess={() => setShowAddToAlbumModal(false)}
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
