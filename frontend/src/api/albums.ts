import { apiClient, EndpointParams } from "./api";
import { Album, ImageMetadata } from "./model";

export const AlbumsAPI = {
  baseEndpoint: "/albums",

  listAlbums: async function () {
    const params: EndpointParams = {
      endpoint: this.baseEndpoint,
      includeCookies: true,
      options: {
        method: "GET",
      },
    };

    return await apiClient<Album[]>(params);
  },

  createAlbum: async function (name: string, description: string = "") {
    const params: EndpointParams = {
      endpoint: this.baseEndpoint,
      includeCookies: true,
      options: {
        method: "POST",
        body: JSON.stringify({ name, description }),
      },
    };

    return await apiClient<Album>(params);
  },

  getAlbum: async function (albumId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}`,
      includeCookies: true,
      options: {
        method: "GET",
      },
    };

    return await apiClient<Album>(params);
  },

  updateAlbum: async function (
    albumId: string,
    name: string,
    description: string = "",
  ) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}`,
      includeCookies: true,
      options: {
        method: "PUT",
        body: JSON.stringify({ name, description }),
      },
    };

    return await apiClient<{ message: string }>(params);
  },

  deleteAlbum: async function (albumId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}`,
      includeCookies: true,
      options: {
        method: "DELETE",
      },
    };

    return await apiClient<{ message: string }>(params);
  },

  getAlbumImages: async function (albumId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}/images`,
      includeCookies: true,
      options: {
        method: "GET",
      },
    };

    return await apiClient<ImageMetadata[]>(params);
  },

  addAlbumImage: async function (albumId: string, imageId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}/images`,
      includeCookies: true,
      options: {
        method: "POST",
        body: JSON.stringify({ image_id: imageId }),
      },
    };

    return await apiClient<{ message: string }>(params);
  },

  deleteAlbumImage: async function (albumId: string, imageId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${albumId}/images/${imageId}`,
      includeCookies: true,
      options: {
        method: "DELETE",
      },
    };

    return await apiClient<{ message: string }>(params);
  },
};
