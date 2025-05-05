import { apiClient, EndpointParams } from "./api";
import { ImageMetadata } from "./model";

export const ImagesAPI = {
  baseEndpoint: "/images",

  listImages: function () {
    const params: EndpointParams = {
      endpoint: this.baseEndpoint,
      requiresAuth: true,
      options: {
        method: "GET",
      },
    };

    return apiClient<ImageMetadata[]>(params);
  },

  uploadImage: function (file: File) {
    // check if correct headers are being set
    const formData = new FormData();
    formData.append("file", file);

    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/upload`,
      requiresAuth: true,
      options: {
        method: "POST",
        body: formData,
      },
    };

    return apiClient<ImageMetadata>(params);
  },

  getImageMetadata: function (imageId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${imageId}`,
      requiresAuth: true,
      options: {
        method: "GET",
      },
    };

    return apiClient<ImageMetadata>(params);
  },

  deleteImage: function (imageId: string) {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/${imageId}`,
      requiresAuth: true,
      options: {
        method: "DELETE",
      },
    };

    return apiClient<{ message: string }>(params);
  },

  getImageURL: function (imageId: string) {
    return `/api${this.baseEndpoint}/${imageId}/download`;
  },
};
