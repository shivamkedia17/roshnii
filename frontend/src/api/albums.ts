import { apiClient } from "./api";
import {
  AlbumInfo,
  CreateAlbumRequest,
  UpdateAlbumRequest,
  PhotoInfo,
} from "@/types";

export const albumsAPI = {
  // Get all albums
  getAlbums: () => apiClient<AlbumInfo[]>("/albums"),

  // Get a single album
  getAlbum: (id: number) => apiClient<AlbumInfo>(`/albums/${id}`),

  // Create a new album
  createAlbum: (data: CreateAlbumRequest) =>
    apiClient<AlbumInfo>("/albums", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  // Update an album
  updateAlbum: (id: number, data: UpdateAlbumRequest) =>
    apiClient<{ message: string }>(`/albums/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),

  // Delete an album
  deleteAlbum: (id: number) =>
    apiClient<{ message: string }>(`/albums/${id}`, {
      method: "DELETE",
    }),

  // Get images in an album
  getAlbumImages: (id: number) =>
    apiClient<PhotoInfo[]>(`/albums/${id}/images`),

  // Add image to album
  addImageToAlbum: (albumId: number, imageId: string) =>
    apiClient<{ message: string }>(`/albums/${albumId}/images`, {
      method: "POST",
      body: JSON.stringify({ image_id: imageId }),
    }),

  // Remove image from album
  removeImageFromAlbum: (albumId: number, imageId: string) =>
    apiClient<{ message: string }>(`/albums/${albumId}/images/${imageId}`, {
      method: "DELETE",
    }),
};
