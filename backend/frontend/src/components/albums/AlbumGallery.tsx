import { useState } from "react";
import { useListAlbums, useDeleteAlbum } from "@/hooks/useAlbums";
import { AlbumItem } from "./AlbumItem";
import { NewAlbumModal } from "./NewAlbumModal";
import "@/css/Albums.css";

export interface AlbumGalleryProps {
  onSelectAlbum: (albumId: string) => void;
}

export function AlbumGallery({ onSelectAlbum }: AlbumGalleryProps) {
  const [showNewAlbumModal, setShowNewAlbumModal] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState<string | null>(null);

  // Fetch albums
  const { data: albums, isLoading, error, refetch } = useListAlbums();

  const deleteAlbum = useDeleteAlbum();

  const handleDeleteAlbum = (albumId: string) => {
    setConfirmDelete(albumId);
  };

  const confirmDeleteAlbum = () => {
    if (!confirmDelete) return;

    deleteAlbum(confirmDelete);
    setConfirmDelete(null);
  };

  if (isLoading) {
    return <div className="loading">Loading albums...</div>;
  }

  if (error) {
    return (
      <div className="error-container">
        <p>
          Error loading albums:{" "}
          {error instanceof Error ? error.message : "Unknown error"}
        </p>
        <button onClick={() => refetch()} className="retry-button">
          Retry
        </button>
      </div>
    );
  }

  return (
    <div className="album-container">
      <div className="album-header">
        <h2>Your Albums</h2>
        <button
          className="create-album-button"
          onClick={() => setShowNewAlbumModal(true)}
        >
          Create New Album
        </button>
      </div>

      {!albums || albums.length === 0 ? (
        <div className="no-albums">
          <p>You don't have any albums yet. Create one to get started!</p>
        </div>
      ) : (
        <div className="albums-grid">
          {albums.map((album) => (
            <AlbumItem
              key={album.id}
              album={album}
              onSelect={onSelectAlbum}
              onDelete={handleDeleteAlbum}
            />
          ))}
        </div>
      )}

      {/* New Album Modal */}
      {showNewAlbumModal && (
        <NewAlbumModal
          onClose={() => setShowNewAlbumModal(false)}
          onSuccess={() => refetch()}
        />
      )}

      {/* Delete Confirmation Dialog */}
      {confirmDelete && (
        <div className="delete-confirm-dialog">
          <h4>Delete Album?</h4>
          <p>This action cannot be undone.</p>

          {/* {deleteAlbum.error && (
            <div className="error">
              {deleteAlbum.error instanceof Error
                ? deleteAlbum.error.message
                : "Failed to delete album"}
            </div>
          )} */}

          <div className="confirm-buttons">
            <button
              onClick={() => setConfirmDelete(null)}
              className="cancel-button"
              // disabled={deleteAlbum.isPending}
            >
              Cancel
            </button>
            <button
              onClick={confirmDeleteAlbum}
              className="confirm-delete-button"
              // disabled={deleteAlbum.isPending}
            >
              {/* {deleteAlbum.isPending ? "Deleting..." : "Confirm Delete"} */}
              {"Confirm Delete"}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
