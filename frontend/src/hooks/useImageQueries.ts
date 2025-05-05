// src/hooks/useImages.ts
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ImagesAPI } from "@/api/images";

// Query keys
export const imageKeys = {
  all: ["images"] as const,
  lists: () => [...imageKeys.all, "list"] as const,
  list: (filters: string) => [...imageKeys.lists(), { filters }] as const,
  details: () => [...imageKeys.all, "detail"] as const,
  detail: (id: string) => [...imageKeys.details(), id] as const,
};

// Hook for fetching all images
export function useImages() {
  return useQuery({
    queryKey: imageKeys.lists(),
    queryFn: () => ImagesAPI.listImages(),
  });
}

// Hook for fetching a single image's metadata
export function useImageMetadata(imageId: string) {
  return useQuery({
    queryKey: imageKeys.detail(imageId),
    queryFn: () => ImagesAPI.getImageMetadata(imageId),
    enabled: !!imageId,
  });
}

// Hook for uploading an image
export function useUploadImage() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (file: File) => ImagesAPI.uploadImage(file),
    onSuccess: () => {
      // Invalidate the images list cache to trigger a refetch
      queryClient.invalidateQueries({ queryKey: imageKeys.lists() });
    },
  });
}

// Hook for deleting an image
export function useDeleteImage() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (imageId: string) => ImagesAPI.deleteImage(imageId),
    onSuccess: (_, imageId) => {
      // Remove the image from cache and invalidate the images list
      queryClient.removeQueries({ queryKey: imageKeys.detail(imageId) });
      queryClient.invalidateQueries({ queryKey: imageKeys.lists() });
    },
  });
}

// Utility function to get image URL (not a React Query hook)
export function getImageURL(imageId: string) {
  return ImagesAPI.getImageURL(imageId);
}
