import axiosInstance from "./api";
import { Album, ImageMetadata } from "./model";

export const AlbumsAPI = {
  listAlbums: async function () {
    const response = await axiosInstance.get("/albums");
    return response.data as Album[];
  },

  createAlbum: async function (name: string, description: string = "") {
    const response = await axiosInstance.post("/albums", { name, description });
    return response.data as Album;
  },

  getAlbum: async function (albumId: string) {
    const response = await axiosInstance.get(`/albums/${albumId}`);
    return response.data as Album;
  },

  updateAlbum: async function (
    albumId: string,
    name: string,
    description: string = "",
  ) {
    const response = await axiosInstance.put(`/albums/${albumId}`, {
      name,
      description,
    });
    return response.data;
  },

  deleteAlbum: async function (albumId: string) {
    const response = await axiosInstance.delete(`/albums/${albumId}`);
    return response.data;
  },

  getAlbumImages: async function (albumId: string) {
    const response = await axiosInstance.get(`/albums/${albumId}/images`);
    return response.data as ImageMetadata[];
  },

  addAlbumImage: async function (albumId: string, imageId: string) {
    const response = await axiosInstance.post(`/albums/${albumId}/images`, {
      image_id: imageId,
    });
    return response.data;
  },

  deleteAlbumImage: async function (albumId: string, imageId: string) {
    const response = await axiosInstance.delete(
      `/albums/${albumId}/images/${imageId}`,
    );
    return response.data;
  },
};
