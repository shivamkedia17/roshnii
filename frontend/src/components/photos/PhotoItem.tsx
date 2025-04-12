import "../../css/Gallery.css";

interface PhotoItemProps {
  photo: {
    id: string;
    thumbnailUrl: string;
    title?: string;
  };
  onClick: () => void;
}

export function PhotoItem({ photo, onClick }: PhotoItemProps) {
  return (
    <div className="photo-item" onClick={onClick}>
      <img
        src={photo.thumbnailUrl}
        alt={photo.title || "Photo"}
        className="photo-thumbnail"
      />
    </div>
  );
}
