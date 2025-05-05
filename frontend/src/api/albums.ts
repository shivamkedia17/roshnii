import { apiClient, EndpointParams } from "./api";
import { Album, ImageMetadata } from "./model";

export const AlbumsAPI = {
  baseEndpoint: "/albums",

  listAlbums: function () {
    const params: EndpointParams = {
      endpoint: this.baseEndpoint,
      requiresAuth: true,
      options: {
        method: "GET",
      },
    };

    return apiClient<Album[]>(params);
  },

  createAlbum: function (name: string, description: string = "") {
    const params: EndpointParams = {
      endpoint: this.baseEndpoint,
      requiresAuth: true,
      options: {
        method: "POST",
        body: JSON.stringify({ name, description }),
      },
    };

    return apiClient<Album>(params);
  },

  getAlbum: function (albumId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}`,
      requiresAuth: true,
      options: {
        method: "GET",
      },
    };

    return apiClient<Album>(params);
  },

  updateAlbum: function (
    albumId: string,
    name: string,
    description: string = "",
  ) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}`,
      requiresAuth: true,
      options: {
        method: "PUT",
        body: JSON.stringify({ name, description }),
      },
    };

    return apiClient<{ message: string }>(params);
  },

  deleteAlbum: function (albumId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}`,
      requiresAuth: true,
      options: {
        method: "DELETE",
      },
    };

    return apiClient<{ message: string }>(params);
  },

  getAlbumImages: function (albumId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}/images`,
      requiresAuth: true,
      options: {
        method: "GET",
      },
    };

    return apiClient<ImageMetadata[]>(params);
  },

  addAlbumImage: function (albumId: string, imageId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}/images`,
      requiresAuth: true,
      options: {
        method: "POST",
        body: JSON.stringify({ image_id: imageId }),
      },
    };

    return apiClient<{ message: string }>(params);
  },

  deleteAlbumImage: function (albumId: string, imageId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}/images/${imageId}`,
      requiresAuth: true,
      options: {
        method: "DELETE",
      },
    };

    return apiClient<{ message: string }>(params);
  },
};
