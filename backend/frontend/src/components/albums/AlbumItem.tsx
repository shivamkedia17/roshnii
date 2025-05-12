import { Album } from "@/api/model";
import "@/css/Albums.css";

interface AlbumItemProps {
  album: Album;
  onSelect: (albumId: string) => void;
  onDelete: (albumId: string) => void;
}

export function AlbumItem({ album, onSelect, onDelete }: AlbumItemProps) {
  return (
    <div className="album-card">
      <div className="album-thumbnail" onClick={() => onSelect(album.id)}>
        {/* Placeholder thumbnail */}
        <div className="placeholder-thumb">{album.name.charAt(0)}</div>
      </div>
      <div className="album-info">
        <h3>{album.name}</h3>
        {album.description && (
          <p className="album-description">{album.description} </p>
        )}
        <p className="album-date">
          Updated {new Date(album.updated_at).toLocaleDateString()}
        </p>
        <div className="album-actions">
          <button
            className="view-album-button"
            onClick={() => onSelect(album.id)}
          >
            View
          </button>
          <button
            className="delete-album-button"
            onClick={(e) => {
              e.stopPropagation();
              onDelete(album.id);
            }}
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  );
}
