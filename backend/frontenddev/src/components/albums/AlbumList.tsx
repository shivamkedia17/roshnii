import { useState } from "react";
import {
  useAlbums,
  useCreateAlbum,
  useDeleteAlbum,
} from "@/hooks/useAlbumQueries";
import { CreateAlbumRequest } from "@/types";
import "@/css/Albums.css";

export type AlbumListProps = {
  onSelectAlbum: (id: string) => void;
};

export function AlbumList({ onSelectAlbum }: AlbumListProps) {
  // State for new album form
  const [showNewAlbumForm, setShowNewAlbumForm] = useState(false);
  const [newAlbumName, setNewAlbumName] = useState("");
  const [newAlbumDescription, setNewAlbumDescription] = useState("");

  // State for album being deleted
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  // Use album hooks
  const { data: albums, isLoading, error, refetch } = useAlbums();
  const createAlbumMutation = useCreateAlbum();
  const deleteAlbumMutation = useDeleteAlbum();

  // Handler to create a new album
  const handleCreateAlbum = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newAlbumName.trim()) return;

    const albumData: CreateAlbumRequest = {
      name: newAlbumName.trim(),
      description: newAlbumDescription.trim() || undefined,
    };

    try {
      await createAlbumMutation.mutateAsync(albumData);
      // Reset form
      setNewAlbumName("");
      setNewAlbumDescription("");
      setShowNewAlbumForm(false);
    } catch (err) {
      console.error("Failed to create album:", err);
    }
  };

  // Handler to delete an album
  const handleDeleteAlbum = async (albumId: number) => {
    try {
      await deleteAlbumMutation.mutateAsync(albumId);
      setConfirmDelete(null);
    } catch (err) {
      console.error("Failed to delete album:", err);
    }
  };

  if (isLoading) return <div className="loading">Loading albums...</div>;
  if (error)
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

  return (
    <div className="album-container">
      <div className="album-header">
        <h2>Your Albums</h2>
        <button
          className="create-album-button"
          onClick={() => setShowNewAlbumForm(true)}
        >
          Create New Album
        </button>
      </div>

      {/* New Album Form */}
      {showNewAlbumForm && (
        <div className="new-album-form">
          <h3>Create New Album</h3>
          <form onSubmit={handleCreateAlbum}>
            <div className="form-group">
              <label htmlFor="album-name">Album Name *</label>
              <input
                id="album-name"
                type="text"
                value={newAlbumName}
                onChange={(e) => setNewAlbumName(e.target.value)}
                required
                placeholder="Enter album name"
              />
            </div>
            <div className="form-group">
              <label htmlFor="album-description">Description (Optional)</label>
              <textarea
                id="album-description"
                value={newAlbumDescription}
                onChange={(e) => setNewAlbumDescription(e.target.value)}
                placeholder="Enter album description"
                rows={3}
              />
            </div>
            <div className="form-buttons">
              <button
                type="button"
                className="cancel-button"
                onClick={() => setShowNewAlbumForm(false)}
              >
                Cancel
              </button>
              <button
                type="submit"
                className="create-button"
                disabled={createAlbumMutation.isPending || !newAlbumName.trim()}
              >
                {createAlbumMutation.isPending ? "Creating..." : "Create Album"}
              </button>
            </div>
            {createAlbumMutation.error && (
              <div className="error">
                {createAlbumMutation.error instanceof Error
                  ? createAlbumMutation.error.message
                  : "Failed to create album"}
              </div>
            )}
          </form>
        </div>
      )}

      {(!albums || albums.length === 0) && !showNewAlbumForm ? (
        <div className="no-albums">
          <p>You don't have any albums yet. Create one to get started!</p>
          <button
            className="create-album-button-empty"
            onClick={() => setShowNewAlbumForm(true)}
          >
            Create Your First Album
          </button>
        </div>
      ) : (
        <div className="albums-grid">
          {albums?.map((album) => (
            <div key={album.id} className="album-card">
              <div
                className="album-thumbnail"
                onClick={() => onSelectAlbum(album.id.toString())}
              >
                {/* Placeholder thumbnail */}
                <div className="placeholder-thumb">{album.name.charAt(0)}</div>
              </div>
              <div className="album-info">
                <h3>{album.name}</h3>
                <p className="album-description">{album.description}</p>
                <p className="album-date">
                  Updated {new Date(album.updated_at).toLocaleDateString()}
                </p>
                <div className="album-actions">
                  <button
                    className="view-album-button"
                    onClick={() => onSelectAlbum(album.id.toString())}
                  >
                    View
                  </button>
                  <button
                    className="delete-album-button"
                    onClick={() => setConfirmDelete(album.id)}
                  >
                    Delete
                  </button>
                </div>
              </div>

              {/* Delete confirmation dialog */}
              {confirmDelete === album.id && (
                <div className="delete-confirm-dialog">
                  <h4>Delete Album?</h4>
                  <p>This action cannot be undone.</p>
                  {deleteAlbumMutation.error && (
                    <div className="error">
                      {deleteAlbumMutation.error instanceof Error
                        ? deleteAlbumMutation.error.message
                        : "Failed to delete album"}
                    </div>
                  )}
                  <div className="confirm-buttons">
                    <button
                      onClick={() => setConfirmDelete(null)}
                      className="cancel-button"
                      disabled={deleteAlbumMutation.isPending}
                    >
                      Cancel
                    </button>
                    <button
                      onClick={() => handleDeleteAlbum(album.id)}
                      className="confirm-delete-button"
                      disabled={deleteAlbumMutation.isPending}
                    >
                      {deleteAlbumMutation.isPending
                        ? "Deleting..."
                        : "Confirm Delete"}
                    </button>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
