import { useState } from "react";
import { useListImages } from "@/hooks/useImages";
import { useListAlbumImages } from "@/hooks/useAlbums";
import { PhotoItem } from "./PhotoItem";
import { PhotoModal } from "./PhotoModal";
import "@/css/Gallery.css";
import { ImageID } from "@/api/model";

export type GalleryProps = {
  albumId?: string;
};

// Either an Album's Photos are being shown, or all of them
export function Gallery({ albumId }: GalleryProps) {
  const [selectedImageId, setSelectedImageId] = useState<ImageID | null>(null);

  // Use the appropriate query based on whether we're viewing all images or an album
  const result = albumId ? useListAlbumImages(albumId) : useListImages();

  const { data: images, isLoading, error } = result;

  if (isLoading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner"></div>
        <p>Loading photos...</p>
      </div>
    );
  }

  // FIXME
  if (error) {
    return (
      <div className="error-container">
        <p>
          Error loading photos:{" "}
          {error instanceof Error ? error.message : "Unknown error"}
        </p>
      </div>
    );
  }

  return (
    <div className="photo-gallery">
      {!images || images.length === 0 ? (
        <div className="no-photos">
          <div className="no-photos-icon">📷</div>
          <p>
            No photos found.{" "}
            {albumId
              ? "This album is empty. Add some photos!"
              : "Upload some photos to get started!"}
          </p>
        </div>
      ) : (
        <>
          <div className="gallery-header">
            <h2>{albumId ? "Album Photos" : "All Photos"}</h2>
            <p className="photo-count">{images.length} Photo(s)</p>
          </div>

          <div className="gallery-grid">
            {images.map((image) => (
              <PhotoItem
                key={image.id}
                imageDetails={image}
                onClick={() => setSelectedImageId(image.id)}
              />
            ))}
          </div>
        </>
      )}

      {selectedImageId && (
        <PhotoModal
          imageId={selectedImageId}
          albumId={albumId}
          onClose={() => setSelectedImageId(null)}
        />
      )}
    </div>
  );
}
