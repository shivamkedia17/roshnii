import { createContext, useContext, ReactNode } from "react";
import { usePhotos as useTanStackPhotos } from "@/hooks/usePhotoQueries"; // Hook that actually fetches data

type PhotoContextType = {
  refreshPhotos: () => void;
};

const PhotoContext = createContext<PhotoContextType | undefined>(undefined);

export function PhotoProvider({ children }: { children: ReactNode }) {
  // Using empty parameters to satisfy the hook's signature
  const { refetch } = useTanStackPhotos({});

  const refreshPhotos = () => {
    refetch();
  };

  return (
    <PhotoContext.Provider value={{ refreshPhotos }}>
      {children}
    </PhotoContext.Provider>
  );
}

// This is now a context consumer function, not a data fetching hook
export function usePhotoContext() {
  const context = useContext(PhotoContext);
  if (!context) {
    throw new Error("usePhotoContext must be used within a PhotoProvider");
  }
  return context;
}
