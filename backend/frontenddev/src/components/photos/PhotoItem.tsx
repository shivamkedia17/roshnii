import { useState } from "react";
import "@/css/Gallery.css";
import { PhotoInfo } from "@/types";

interface PhotoItemProps {
  photo: PhotoInfo;
  onClick: () => void;
}

export function PhotoItem({ photo, onClick }: PhotoItemProps) {
  const [isLoading, setIsLoading] = useState(true);
  const [hasError, setHasError] = useState(false);

  return (
    <div className="photo-item" onClick={onClick}>
      {isLoading && <div className="photo-loading"></div>}

      {hasError && (
        <div className="photo-error">
          <span>⚠️</span>
        </div>
      )}

      <img
        src={photo.thumbnailUrl}
        alt={photo.filename || "Photo"}
        className={`photo-thumbnail ${isLoading ? "loading" : ""} ${hasError ? "error" : ""}`}
        onLoad={() => setIsLoading(false)}
        onError={() => {
          setIsLoading(false);
          setHasError(true);
        }}
      />
    </div>
  );
}
