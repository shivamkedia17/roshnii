// UI Hooks to fetch Albums and Metadata

import { AlbumsAPI } from "@/api/albums";
import { AlbumID, ImageID, ImageMetadata } from "@/api/model";
import {
  useMutation,
  useQuery,
  useQueryClient,
  keepPreviousData,
} from "@tanstack/react-query";

// Query keys
export const albumKeys = {
  all: ["albums"] as const,
  lists: () => [...albumKeys.all, "list"] as const,
  list: (filters: string) => [...albumKeys.lists(), { filters }] as const,
  details: () => [...albumKeys.all, "detail"] as const,
  detail: (id: AlbumID) => [...albumKeys.details(), id] as const,
  images: (id: AlbumID) => [...albumKeys.detail(id), "images"] as const,
};

// Hook to list all albums
export function useListAlbums() {
  return useQuery({
    queryKey: albumKeys.lists(),
    queryFn: AlbumsAPI.listAlbums,
    placeholderData: keepPreviousData,
    throwOnError: true,
  });
}

// Hook to get a specific album
export function useListAlbum(albumId: AlbumID) {
  return useQuery({
    queryKey: albumKeys.detail(albumId),
    queryFn: () => AlbumsAPI.listAlbum(albumId),
    placeholderData: keepPreviousData,
    throwOnError: true,
  });
}

// Hook to get images in an album
export function useListAlbumImages(albumId: AlbumID) {
  return useQuery({
    queryKey: albumKeys.images(albumId),
    queryFn: () => AlbumsAPI.listAlbumImages(albumId),
    placeholderData: keepPreviousData,
    throwOnError: true,
  });
}

// Mutation hook to create an album
export function useCreateAlbum() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: ({
      name,
      description,
    }: {
      name: string;
      description?: string;
    }) => AlbumsAPI.createAlbum(name, description),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: albumKeys.lists() });
    },
  });

  return mutation.mutate;
}

// Mutation hook to update an album
export function useUpdateAlbum() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: ({
      albumId,
      name,
      description,
    }: {
      albumId: AlbumID;
      name: string;
      description?: string;
    }) => AlbumsAPI.updateAlbum(albumId, name, description),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: albumKeys.detail(variables.albumId),
      });
      queryClient.invalidateQueries({ queryKey: albumKeys.lists() });
    },
  });

  return mutation.mutate;
}

// Mutation hook to delete an album
export function useDeleteAlbum() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (albumId: AlbumID) => AlbumsAPI.deleteAlbum(albumId),
    onSuccess: (_, albumId) => {
      queryClient.invalidateQueries({ queryKey: albumKeys.lists() });
      // Remove the specific album from cache
      queryClient.removeQueries({ queryKey: albumKeys.detail(albumId) });
    },
  });

  return mutation.mutate;
}

// Mutation hook to add an image to an album
export function useAddImageToAlbum() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: ({
      albumId,
      imageId,
    }: {
      albumId: AlbumID;
      imageId: ImageID;
    }) => AlbumsAPI.addAlbumImage(albumId, imageId),
    onSuccess: (_, variables) => {
      // Invalidate the album images query to refresh the list
      queryClient.invalidateQueries({
        queryKey: albumKeys.images(variables.albumId),
      });

      // Optionally update the cache directly for a more responsive UI
      queryClient.setQueryData<ImageMetadata[]>(
        albumKeys.images(variables.albumId),
        (oldData) => {
          // If we have the current image list cached
          if (oldData) {
            // Check if image is already in the album
            const imageExists = oldData.some(
              (img) => img.id === variables.imageId,
            );
            if (!imageExists) {
              // We could fetch the image metadata and add it,
              // but for simplicity we'll just invalidate the query
              return undefined; // Force a refetch
            }
            return oldData;
          }
          return undefined; // Force a refetch
        },
      );
    },
  });

  return mutation.mutate;
}

// Mutation hook to remove an image from an album
export function useRemoveImageFromAlbum() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: ({
      albumId,
      imageId,
    }: {
      albumId: AlbumID;
      imageId: ImageID;
    }) => AlbumsAPI.deleteAlbumImage(albumId, imageId),
    onSuccess: (_, variables) => {
      // Invalidate the album images query
      queryClient.invalidateQueries({
        queryKey: albumKeys.images(variables.albumId),
      });

      // Optionally update the cache directly for immediate UI feedback
      queryClient.setQueryData<ImageMetadata[]>(
        albumKeys.images(variables.albumId),
        (oldData) => {
          if (oldData) {
            // Remove the image from the cached data
            return oldData.filter((img) => img.id !== variables.imageId);
          }
          return undefined; // Force a refetch
        },
      );
    },
  });

  return mutation.mutate;
}
