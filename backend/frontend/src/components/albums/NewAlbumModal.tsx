import { useState } from "react";
import { useCreateAlbum } from "@/hooks/useAlbums";
import "@/css/AddToAlbumModal.css";

interface NewAlbumModalProps {
  onClose: () => void;
  onSuccess?: () => void;
}

export function NewAlbumModal({ onClose, onSuccess }: NewAlbumModalProps) {
  const [albumName, setAlbumName] = useState("");
  const [albumDescription, setAlbumDescription] = useState("");

  const createAlbum = useCreateAlbum();
  const handleCreateAlbum = () => {
    if (!albumName.trim()) return;

    createAlbum(
      {
        name: albumName.trim(),
        description: albumDescription.trim() || undefined,
      },
      {
        onSuccess: () => {
          if (onSuccess) {
            onSuccess();
          }
          onClose();
        },
      },
    );
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="add-to-album-modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3>Create New Album</h3>
          <button className="close-button" onClick={onClose}>
            ×
          </button>
        </div>

        <div className="modal-body">
          <form
            onSubmit={(e) => {
              e.preventDefault();
              handleCreateAlbum();
            }}
          >
            <div className="form-group">
              <label htmlFor="album-name">Album Name *</label>
              <input
                id="album-name"
                type="text"
                value={albumName}
                onChange={(e) => setAlbumName(e.target.value)}
                required
                placeholder="Enter album name"
                autoFocus
              />
            </div>
            <div className="form-group">
              <label htmlFor="album-description">Description (Optional)</label>
              <textarea
                id="album-description"
                value={albumDescription}
                onChange={(e) => setAlbumDescription(e.target.value)}
                placeholder="Enter album description"
                rows={3}
              />
            </div>
          </form>
        </div>

        <div className="modal-footer">
          <button className="cancel-button" onClick={onClose}>
            Cancel
          </button>
          <button
            className="add-button"
            onClick={handleCreateAlbum}
            disabled={!albumName.trim()}
          >
            {/* {createAlbum.isPending ? "Creating..." : "Create Album"} */}
            {"Create Album"}
          </button>
        </div>
      </div>
    </div>
  );
}
