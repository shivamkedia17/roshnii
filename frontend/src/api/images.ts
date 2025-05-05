import { apiClient, EndpointParams } from "./api";
import { ImageMetadata } from "./model";

export const ImagesAPI = {
  baseEndpoint: "/images",

  listImages: async function () {
    const params: EndpointParams = {
      endpoint: this.baseEndpoint,
      includeCookies: true,
      options: {
        method: "GET",
      },
    };

    return await apiClient<ImageMetadata[]>(params);
  },

  uploadImage: async function (file: File) {
    // check if correct headers are being set
    const formData = new FormData();
    formData.append("file", file);

    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/upload`,
      includeCookies: true,
      options: {
        method: "POST",
        body: formData,
      },
    };

    return await apiClient<ImageMetadata>(params);
  },

  getImageMetadata: async function (imageId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${imageId}`,
      includeCookies: true,
      options: {
        method: "GET",
      },
    };

    return await apiClient<ImageMetadata>(params);
  },

  deleteImage: async function (imageId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${imageId}`,
      includeCookies: true,
      options: {
        method: "DELETE",
      },
    };

    return await apiClient<{ message: string }>(params);
  },

  loadImage: async function (imageId: string) {
    const url = `/api${this.baseEndpoint}/${imageId}/download`;

    const response = await fetch(url, {
      method: "GET",
      credentials: "include",
    });

    if (!response.ok) {
      throw new Error(
        `Failed to load image: ${response.status} ${response.statusText}`,
      );
    }

    return await response.blob();
  },
};
