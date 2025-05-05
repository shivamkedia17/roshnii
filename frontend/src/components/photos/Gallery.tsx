import { useState } from "react";
import { PhotoItem } from "./PhotoItem";
import { PhotoModal } from "./PhotoModal";
import { useImages } from "@/hooks/useImageQueries";
import "@/css/Gallery.css";

export type GalleryProps = {
  searchQuery?: string;
  albumId?: string;
};

export function Gallery({
  searchQuery = "",
  albumId = undefined,
}: GalleryProps) {
  const [selectedPhoto, setSelectedPhoto] = useState<string | null>(null);
  const {
    data: photos,
    isLoading,
    error,
    refetch,
  } = useImages({ searchQuery, albumId });

  // Handler for retrying if there's an error
  const handleRetry = () => {
    refetch();
  };

  if (isLoading)
    return (
      <div className="loading-container">
        <div className="loading-spinner"></div>
        <p>Loading photos...</p>
      </div>
    );

  if (error)
    return (
      <div className="error-container">
        <p>
          Error loading photos:{" "}
          {error instanceof Error ? error.message : "Unknown error"}
        </p>
        <button onClick={handleRetry} className="retry-button">
          Retry
        </button>
      </div>
    );

  return (
    <div className="photo-gallery">
      {!photos || photos.length === 0 ? (
        <div className="no-photos">
          <div className="no-photos-icon">📷</div>
          <p>
            No photos found.{" "}
            {searchQuery
              ? "Try a different search term."
              : albumId
                ? "This album is empty. Add some photos!"
                : "Upload some photos to get started!"}
          </p>
          {!albumId && !searchQuery && (
            <button
              onClick={() => (window.location.hash = "#upload")}
              className="upload-now-button"
            >
              Upload Now
            </button>
          )}
        </div>
      ) : (
        <>
          <div className="gallery-header">
            <h2>
              {albumId
                ? "Album Photos"
                : searchQuery
                  ? `Search Results: "${searchQuery}"`
                  : "Your Photos"}
            </h2>
            <p className="photo-count">{photos.length} photos</p>
          </div>

          <div className="gallery-grid">
            {photos.map((photo) => (
              <PhotoItem
                key={photo.id}
                photo={photo}
                onClick={() => setSelectedPhoto(photo.id)}
              />
            ))}
          </div>
        </>
      )}

      {selectedPhoto && (
        <PhotoModal
          photoId={selectedPhoto}
          onClose={() => setSelectedPhoto(null)}
          albumContext={undefined}
        />
      )}
    </div>
  );
}
