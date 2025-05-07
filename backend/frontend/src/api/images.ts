import axiosInstance from "./api";
import { ImageMetadata } from "./model";

export const ImagesAPI = {
  listImages: async function () {
    const response = await axiosInstance.get("/images");
    return response.data as ImageMetadata[];
  },

  uploadImage: async function (file: File) {
    const formData = new FormData();
    formData.append("file", file);

    const response = await axiosInstance.post("/images/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });

    return response.data as ImageMetadata;
  },

  getImageMetadata: async function (imageId: string) {
    const response = await axiosInstance.get(`/images/${imageId}`);
    return response.data as ImageMetadata;
  },

  deleteImage: async function (imageId: string) {
    const response = await axiosInstance.delete(`/images/${imageId}`);
    return response.data;
  },

  loadImage: async function (imageId: string) {
    const response = await axiosInstance.get(`/images/${imageId}/download`, {
      responseType: "blob",
    });

    return response.data;
  },
};
