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
    mutationFn: (formData: FormData) => {
      // Make sure we're using the correct API endpoint
      const token = localStorage.getItem("auth_token");
      const headers: HeadersInit = {};

      // Add Authorization header if token exists (for dev mode)
      if (token) {
        headers["Authorization"] = `Bearer ${token}`;
      }

      return fetch(`/api/upload`, {
        method: "POST",
        credentials: "include", // For cookies
        headers,
        body: formData, // Don't set Content-Type header with FormData
      }).then((res) => {
        if (!res.ok) {
          return res.json().then((err) => {
            throw new Error(err.error || `Upload failed: ${res.status}`);
          });
        }
        return res.json();
      });
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
    mutationFn: (photoId: string) => {
      const token = localStorage.getItem("auth_token");
      const headers: HeadersInit = {
        "Content-Type": "application/json",
      };

      // Add Authorization header if token exists
      if (token) {
        headers["Authorization"] = `Bearer ${token}`;
      }

      return fetch(`/api/image/${photoId}`, {
        method: "DELETE",
        credentials: "include", // For cookies
        headers,
      }).then((res) => {
        if (!res.ok) {
          return res
            .json()
            .then((err) => {
              throw new Error(err.error || `Delete failed: ${res.status}`);
            })
            .catch(() => {
              // If JSON parsing fails
              throw new Error(`Failed to delete photo: ${res.status}`);
            });
        }
        return res.json();
      });
    },
    onSuccess: () => {
      // Invalidate and refetch photos list when a photo is deleted
      queryClient.invalidateQueries({ queryKey: photoKeys.lists() });
    },
  });
}
