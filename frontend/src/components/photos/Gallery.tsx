import { useEffect, useState } from "react";
import { PhotoItem } from "./PhotoItem";
import { PhotoModal } from "./PhotoModal";
import { usePhotos } from "@/context/PhotoContext";
import { GalleryProps } from "@/types";
import "@/css/Gallery.css";

export function Gallery({
  searchQuery = "",
  albumId = undefined,
}: GalleryProps) {
  const [selectedPhoto, setSelectedPhoto] = useState<string | null>(null);
  const { photos, loading, error } = usePhotos({ searchQuery, albumId });

  if (loading) return <div className="loading">Loading photos...</div>;
  if (error) return <div className="error">Error loading photos: {error}</div>;

  return (
    <div className="photo-gallery">
      {photos.length === 0 ? (
        <div className="no-photos">
          <p>
            No photos found.{" "}
            {searchQuery
              ? "Try a different search term."
              : "Upload some photos to get started!"}
          </p>
        </div>
      ) : (
        <div className="gallery-grid">
          {photos.map((photo) => (
            <PhotoItem
              key={photo.id}
              photo={photo}
              onClick={() => setSelectedPhoto(photo.id)}
            />
          ))}
        </div>
      )}

      {selectedPhoto && (
        <PhotoModal
          photoId={selectedPhoto}
          onClose={() => setSelectedPhoto(null)}
        />
      )}
    </div>
  );
}
