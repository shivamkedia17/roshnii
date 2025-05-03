import { useState } from "react";
import { useAlbums, useAddImageToAlbum } from "@/hooks/useAlbumQueries";
import "@/css/AddToAlbumModal.css";

type AddToAlbumModalProps = {
  imageId: string;
  onClose: () => void;
  onSuccess?: () => void;
};

export function AddToAlbumModal({
  imageId,
  onClose,
  onSuccess,
}: AddToAlbumModalProps) {
  const [selectedAlbumId, setSelectedAlbumId] = useState<number | null>(null);

  // Get albums
  const { data: albums, isLoading, error } = useAlbums();

  // Add to album mutation
  const addToAlbumMutation = useAddImageToAlbum();

  // Handle adding to album
  const handleAddToAlbum = async () => {
    if (!selectedAlbumId) return;

    try {
      await addToAlbumMutation.mutateAsync({
        albumId: selectedAlbumId,
        imageId: imageId,
      });

      if (onSuccess) {
        onSuccess();
      }
      onClose();
    } catch (err) {
      console.error("Failed to add to album:", err);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="add-to-album-modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Add to Album</h3>
          <button className="close-button" onClick={onClose}>
            ×
          </button>
        </div>

        <div className="modal-body">
          {isLoading ? (
            <div className="loading">Loading albums...</div>
          ) : error ? (
            <div className="error">
              {error instanceof Error ? error.message : "Failed to load albums"}
            </div>
          ) : !albums || albums.length === 0 ? (
            <div className="no-albums">
              <p>You don't have any albums yet.</p>
              <button
                className="create-album-button"
                onClick={() => {
                  // Close this modal and navigate to albums view
                  onClose();
                  window.location.hash = "#albums";
                }}
              >
                Create an Album
              </button>
            </div>
          ) : (
            <>
              <p className="instruction">
                Select an album to add this photo to:
              </p>
              <div className="album-list">
                {albums.map((album) => (
                  <div
                    key={album.id}
                    className={`album-option ${selectedAlbumId === album.id ? "selected" : ""}`}
                    onClick={() => setSelectedAlbumId(album.id)}
                  >
                    <div className="album-icon">{album.name.charAt(0)}</div>
                    <div className="album-details">
                      <h4>{album.name}</h4>
                      {album.description && (
                        <p className="album-description">{album.description}</p>
                      )}
                    </div>
                    {selectedAlbumId === album.id && (
                      <div className="selected-indicator">✓</div>
                    )}
                  </div>
                ))}
              </div>
            </>
          )}

          {addToAlbumMutation.error && (
            <div className="error">
              {addToAlbumMutation.error instanceof Error
                ? addToAlbumMutation.error.message
                : "Failed to add to album"}
            </div>
          )}
        </div>

        <div className="modal-footer">
          <button className="cancel-button" onClick={onClose}>
            Cancel
          </button>
          <button
            className="add-button"
            onClick={handleAddToAlbum}
            disabled={
              addToAlbumMutation.isPending ||
              !selectedAlbumId ||
              isLoading ||
              !albums ||
              albums.length === 0
            }
          >
            {addToAlbumMutation.isPending ? "Adding..." : "Add to Album"}
          </button>
        </div>
      </div>
    </div>
  );
}
