// UI Hooks to fetch Images and Metadata

import { ImagesAPI } from "@/api/images";
import { ImageID } from "@/api/model";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import placeholderImagePath from "@/placeholders/mountain.jpeg";
import { albumKeys } from "./useAlbums";

let placeholderBlob: Blob;

fetch(placeholderImagePath)
  .then((response) => response.blob())
  .then((blob) => {
    placeholderBlob = blob;
  });

export const imageKeys = {
  all: ["image", "metadata"] as const,
  imageMeta: (id: ImageID) => ["image", "metadata", id] as const,
  imageData: (id: ImageID) => ["image", "data", id] as const,
};

export function useListImages() {
  return useQuery({
    queryKey: imageKeys.all,
    queryFn: ImagesAPI.getImageMetadataAll,
    throwOnError: true,
  });
}

export function useListImage(imageId: ImageID) {
  return useQuery({
    queryKey: imageKeys.imageMeta(imageId),
    queryFn: () => ImagesAPI.getImageMetadata(imageId),
    throwOnError: true,
  });
}

export function useGetImage(imageId: ImageID) {
  return useQuery({
    queryKey: imageKeys.imageData(imageId),
    queryFn: () => ImagesAPI.loadImage(imageId),
    refetchOnWindowFocus: false,
    placeholderData: placeholderBlob,
  });
}

// Mutation hooks

// Hook for uploading a new image
export function useUploadImage() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (file: File) => ImagesAPI.uploadImage(file),
    onSuccess: () => {
      // After successful upload, invalidate the image list query to refresh
      queryClient.invalidateQueries({ queryKey: imageKeys.all });
    },
    onError: (error) => {
      console.error("Error uploading image:", error);
    },
  });

  return mutation.mutate;
}

// Hook for deleting an image
export function useDeleteImage() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (imageId: ImageID) => ImagesAPI.deleteImage(imageId),
    onSuccess: (_, imageId) => {
      // After deleting, invalidate the image list
      queryClient.invalidateQueries({ queryKey: imageKeys.all });
      queryClient.invalidateQueries({ queryKey: albumKeys.all });

      // Also remove the specific image queries from cache
      queryClient.removeQueries({ queryKey: imageKeys.imageMeta(imageId) });
      queryClient.removeQueries({ queryKey: imageKeys.imageData(imageId) });
    },
    onError: (error) => {
      console.error("Error deleting image:", error);
    },
  });

  return mutation.mutate;
}
