import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { albumsAPI } from "@/api/albums";
import { CreateAlbumRequest, UpdateAlbumRequest } from "@/types";

// Query keys
export const albumKeys = {
  all: ["albums"] as const,
  lists: () => [...albumKeys.all, "list"] as const,
  list: () => [...albumKeys.lists()] as const,
  details: () => [...albumKeys.all, "detail"] as const,
  detail: (id: number) => [...albumKeys.details(), id] as const,
  images: () => [...albumKeys.all, "images"] as const,
  albumImages: (id: number) => [...albumKeys.images(), id] as const,
};

// Get all albums hook
export function useAlbums() {
  return useQuery({
    queryKey: albumKeys.list(),
    queryFn: albumsAPI.getAlbums,
  });
}

// Get single album hook
export function useAlbum(id: number | null) {
  return useQuery({
    queryKey: albumKeys.detail(id || 0),
    queryFn: () => albumsAPI.getAlbum(id || 0),
    enabled: !!id, // Only run the query if we have an ID
  });
}

// Get album images hook
export function useAlbumImages(albumId: number | null) {
  return useQuery({
    queryKey: albumKeys.albumImages(albumId || 0),
    queryFn: () => albumsAPI.getAlbumImages(albumId || 0),
    enabled: !!albumId, // Only run the query if we have an ID
  });
}

// Create album hook
export function useCreateAlbum() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateAlbumRequest) => albumsAPI.createAlbum(data),
    onSuccess: () => {
      // Invalidate albums list to refresh it
      queryClient.invalidateQueries({ queryKey: albumKeys.lists() });
    },
  });
}

// Update album hook
export function useUpdateAlbum() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateAlbumRequest }) =>
      albumsAPI.updateAlbum(id, data),
    onSuccess: (_, variables) => {
      // Invalidate specific album and the albums list
      queryClient.invalidateQueries({
        queryKey: albumKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: albumKeys.lists() });
    },
  });
}

// Delete album hook
export function useDeleteAlbum() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => albumsAPI.deleteAlbum(id),
    onSuccess: () => {
      // Invalidate albums list
      queryClient.invalidateQueries({ queryKey: albumKeys.lists() });
    },
  });
}

// Add image to album hook
export function useAddImageToAlbum() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ albumId, imageId }: { albumId: number; imageId: string }) =>
      albumsAPI.addImageToAlbum(albumId, imageId),
    onSuccess: (_, variables) => {
      // Invalidate album images list
      queryClient.invalidateQueries({
        queryKey: albumKeys.albumImages(variables.albumId),
      });
    },
  });
}

// Remove image from album hook
export function useRemoveImageFromAlbum() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ albumId, imageId }: { albumId: number; imageId: string }) =>
      albumsAPI.removeImageFromAlbum(albumId, imageId),
    onSuccess: (_, variables) => {
      // Invalidate album images list
      queryClient.invalidateQueries({
        queryKey: albumKeys.albumImages(variables.albumId),
      });
    },
  });
}
