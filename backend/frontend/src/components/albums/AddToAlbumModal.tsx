import { useState } from "react";
import { useListAlbums, useAddImageToAlbum } from "@/hooks/useAlbums";
import "@/css/AddToAlbumModal.css";

interface AddToAlbumModalProps {
  imageId: string;
  onClose: () => void;
  onSuccess?: () => void;
}

export function AddToAlbumModal({
  imageId,
  onClose,
  onSuccess,
}: AddToAlbumModalProps) {
  const [selectedAlbumId, setSelectedAlbumId] = useState<string | null>(null);

  const { data: albums, isLoading, error } = useListAlbums();
  const addToAlbum = useAddImageToAlbum();

  const handleAddToAlbum = () => {
    if (!selectedAlbumId) return;

    addToAlbum(
      {
        albumId: selectedAlbumId,
        imageId: imageId,
      },
      {
        onSuccess: () => {
          if (onSuccess) {
            onSuccess();
          }
          onClose();
        },
        onError: (error) => {
          console.error("Failed to add to album:", error);
        },
      },
    );
  };

  return (
    <div className="photo-modal-overlay" onClick={onClose}>
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
              <p>You don't have any albums yet. Please create some.</p>
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
                    className={`album-option ${
                      selectedAlbumId === album.id ? "selected" : ""
                    }`}
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
          {/*
          {addToAlbum.error && (
            <div className="error">
              {addToAlbum.error instanceof Error
                ? addToAlbum.error.message
                : "Failed to add to album"}
            </div>
          )} */}
        </div>

        <div className="modal-footer">
          <button className="cancel-button" onClick={onClose}>
            Cancel
          </button>
          <button
            className="add-button"
            onClick={handleAddToAlbum}
            disabled={
              // addToAlbum.isPending ||
              !selectedAlbumId || isLoading || !albums || albums.length === 0
            }
          >
            {/* {addToAlbum.isPending ? "Adding..." : "Add to Album"} */}
            {"Add to Album"}
          </button>
        </div>
      </div>
    </div>
  );
}
