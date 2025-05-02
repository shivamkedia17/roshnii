import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { photosAPI } from "../services/api";
import { GalleryProps } from "@/types";

type Photo = {
  id: string;
  filename: string;
  content_type: string;
  thumbnailUrl: string; // We'll construct this from the backend URL
  created_at: string;
};

type PhotoContextType = {
  photos: Photo[];
  loading: boolean;
  error: string | null;
  refreshPhotos: () => Promise<void>;
};

const PhotoContext = createContext<PhotoContextType | undefined>(undefined);

export function PhotoProvider({ children }: { children: ReactNode }) {
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchPhotos = async () => {
    try {
      setLoading(true);
      const data = await photosAPI.getPhotos();

      // Transform the data to include thumbnailUrl
      const transformedPhotos = data.map((photo) => ({
        ...photo,
        thumbnailUrl: `/api/image/${photo.id}/download`, // Assuming this endpoint exists
      }));

      setPhotos(transformedPhotos);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load photos");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPhotos();
  }, []);

  return (
    <PhotoContext.Provider
      value={{ photos, loading, error, refreshPhotos: fetchPhotos }}
    >
      {children}
    </PhotoContext.Provider>
  );
}

export function usePhotos({ searchQuery = "", albumId }: GalleryProps) {
  const context = useContext(PhotoContext);
  if (!context) {
    throw new Error("usePhotos must be used within a PhotoProvider");
  }

  const { photos, loading, error } = context;

  // Filter photos based on search query
  const filteredPhotos = photos.filter((photo) =>
    photo.filename.toLowerCase().includes(searchQuery.toLowerCase()),
  );

  return { photos: filteredPhotos, loading, error };
}
