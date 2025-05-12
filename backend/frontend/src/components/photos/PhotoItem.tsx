import { ImageMetadata } from "@/api/model";
import "@/css/Gallery.css";
import { useGetImage } from "@/hooks/useImages";
import { useEffect, useState } from "react";

interface PhotoItemProps {
  imageDetails: ImageMetadata;
  onClick: () => void;
}

export function PhotoItem({ imageDetails, onClick }: PhotoItemProps) {
  const { data: imageBlob, isFetching, isError } = useGetImage(imageDetails.id);
  const [imageBlobURL, setImageBlobURL] = useState<string>("");

  // Set up and clean up object URL when blob changes
  useEffect(() => {
    if (imageBlob) {
      const url = URL.createObjectURL(imageBlob);
      setImageBlobURL(url);

      // Clean up function to revoke the object URL when component unmounts or blob changes
      return () => {
        URL.revokeObjectURL(url);
      };
    }
  }, [imageBlob]);

  return (
    <div className="photo-item" onClick={onClick}>
      {isFetching && <div className="photo-loading"></div>}

      {isError ? (
        <div className="photo-error">
          <span>⚠️</span>
          <p>Failed to load image</p>
        </div>
      ) : (
        imageBlobURL && (
          <img
            src={imageBlobURL}
            alt={"Photo"}
            className={`photo-thumbnail ${isFetching ? "loading" : ""}`}
          />
        )
      )}
    </div>
  );
}
