// src/hooks/usePhotoQueries.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { photosAPI } from "@/api/photos";
import { GalleryProps } from "@/components/photos/Gallery";

// Query keys
export const photoKeys = {
  all: ["photos"] as const,
  lists: () => [...photoKeys.all, "list"] as const,
  list: (filters: { searchQuery?: string; albumId?: string }) =>
    [...photoKeys.lists(), filters] as const,
  details: () => [...photoKeys.all, "detail"] as const,
  detail: (id: string) => [...photoKeys.details(), id] as const,
};

// Get all photos hook
// TODO retrieve, photos and metadata
// TODO add pagination
export function usePhotos({ searchQuery = "", albumId }: GalleryProps = {}) {
  return useQuery({
    queryKey: photoKeys.list({ searchQuery, albumId }),
    queryFn: () => photosAPI.getPhotos(),
    select: (data) => {
      // Transform the data to include thumbnailUrl
      const transformedPhotos = data.map((photo) => ({
        ...photo,
        thumbnailUrl: `/api/image/${photo.id}/download`,
      }));

      // Filter photos based on search query if provided
      if (searchQuery) {
        return transformedPhotos.filter((photo) =>
          photo.filename.toLowerCase().includes(searchQuery.toLowerCase()),
        );
      }

      return transformedPhotos;
    },
  });
}

// Get single photo hook
export function usePhoto(id: string | null) {
  return useQuery({
    queryKey: photoKeys.detail(id || ""),
    queryFn: () => photosAPI.getPhoto(id || ""),
    enabled: !!id, // Only run the query if we have an ID
  });
}

// Upload photo hook
export function useUploadPhoto() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (formData: FormData) => {
      try {
        const response = await fetch(`/api/upload`, {
          method: "POST",
          credentials: "include", // Use cookies for authentication
          body: formData, // Don't set Content-Type header with FormData
        });

        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}));
          throw new Error(
            errorData.error || `Upload failed: ${response.status}`,
          );
        }

        return await response.json();
      } catch (error) {
        console.error("Photo upload failed:", error);
        throw error;
      }
    },
    onSuccess: () => {
      // Invalidate and refetch photos list when photo is uploaded
      queryClient.invalidateQueries({ queryKey: photoKeys.lists() });
    },
  });
}

// Delete photo hook
export function useDeletePhoto() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (photoId: string) => {
      try {
        const response = await fetch(`/api/image/${photoId}`, {
          method: "DELETE",
          credentials: "include", // Use cookies for authentication
          headers: {
            "Content-Type": "application/json",
          },
        });

        if (!response.ok) {
          try {
            const errorData = await response.json();
            throw new Error(
              errorData.error || `Delete failed: ${response.status}`,
            );
          } catch (jsonError) {
            // If JSON parsing fails
            throw new Error(`Failed to delete photo: ${response.status}`);
          }
        }

        return await response.json();
      } catch (error) {
        console.error("Photo deletion failed:", error);
        throw error;
      }
    },
    onSuccess: () => {
      // Invalidate and refetch photos list when a photo is deleted
      queryClient.invalidateQueries({ queryKey: photoKeys.lists() });
    },
  });
}
