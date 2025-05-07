import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AlbumsAPI } from "@/api/albums";

// Query keys
export const albumsAPIQueryKeys = {
  all: ["albums"] as const,
  lists: () => [...albumsAPIQueryKeys.all, "list"] as const,
  list: (filters: string) =>
    [...albumsAPIQueryKeys.lists(), { filters }] as const,
  details: () => [...albumsAPIQueryKeys.all, "detail"] as const,
  detail: (id: string) => [...albumsAPIQueryKeys.details(), id] as const,
  images: (id: string) => [...albumsAPIQueryKeys.detail(id), "images"] as const,
};

// Hooks for data fetching
export function useAlbums() {
  return useQuery({
    queryKey: albumsAPIQueryKeys.lists(),
    queryFn: () => AlbumsAPI.listAlbums(),
  });
}

export function useAlbum(albumId: string) {
  return useQuery({
    queryKey: albumsAPIQueryKeys.detail(albumId),
    queryFn: () => AlbumsAPI.getAlbum(albumId),
    enabled: !!albumId,
  });
}

export function useAlbumImages(albumId: string) {
  return useQuery({
    queryKey: albumsAPIQueryKeys.images(albumId),
    queryFn: () => AlbumsAPI.getAlbumImages(albumId),
    enabled: !!albumId,
  });
}

// Hooks for mutations
export function useCreateAlbum() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      name,
      description,
    }: {
      name: string;
      description: string;
    }) => AlbumsAPI.createAlbum(name, description),
    onSuccess: () => {
      // Invalidate the albums list cache to trigger a refetch
      queryClient.invalidateQueries({ queryKey: albumsAPIQueryKeys.lists() });
    },
  });
}

export function useUpdateAlbum() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      albumId,
      name,
      description,
    }: {
      albumId: string;
      name: string;
      description: string;
    }) => AlbumsAPI.updateAlbum(albumId, name, description),
    onSuccess: (_, variables) => {
      // Invalidate the specific album cache and the albums list
      queryClient.invalidateQueries({
        queryKey: albumsAPIQueryKeys.detail(variables.albumId),
      });
      queryClient.invalidateQueries({ queryKey: albumsAPIQueryKeys.lists() });
    },
  });
}

export function useDeleteAlbum() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (albumId: string) => AlbumsAPI.deleteAlbum(albumId),
    onSuccess: (_, albumId) => {
      // Remove the album from cache and invalidate the albums list
      queryClient.removeQueries({
        queryKey: albumsAPIQueryKeys.detail(albumId),
      });
      queryClient.invalidateQueries({ queryKey: albumsAPIQueryKeys.lists() });
    },
  });
}

export function useAddImageToAlbum() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ albumId, imageId }: { albumId: string; imageId: string }) =>
      AlbumsAPI.addAlbumImage(albumId, imageId),
    onSuccess: (_, variables) => {
      // Invalidate the album images cache
      queryClient.invalidateQueries({
        queryKey: albumsAPIQueryKeys.images(variables.albumId),
      });
    },
  });
}

export function useRemoveImageFromAlbum() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ albumId, imageId }: { albumId: string; imageId: string }) =>
      AlbumsAPI.deleteAlbumImage(albumId, imageId),
    onSuccess: (_, variables) => {
      // Invalidate the album images cache
      queryClient.invalidateQueries({
        queryKey: albumsAPIQueryKeys.images(variables.albumId),
      });
    },
  });
}
