import axiosInstance from "./api";
import { ImageID, ImageMetadata, ServerMessage } from "./model";

export const ImagesAPI = {
  getImageMetadataAll: async function () {
    const response = await axiosInstance.get("/images");
    return response.data as ImageMetadata[];
  },

  getImageMetadata: async function (imageId: ImageID) {
    const response = await axiosInstance.get(`/images/${imageId}`);
    return response.data as ImageMetadata;
  },

  deleteImage: async function (imageId: ImageID) {
    const response = await axiosInstance.delete(`/images/${imageId}`);
    return response.data as ServerMessage;
  },

  loadImage: async function (imageId: ImageID) {
    const response = await axiosInstance.get(`/images/${imageId}/download`, {
      responseType: "blob",
    });

    return response.data as Blob;
  },

  uploadImage: async function (file: File) {
    const formData = new FormData();
    formData.append("file", file);

    // TODO? check for valid filetype?

    const response = await axiosInstance.post("/images/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });

    return response.data as ImageMetadata;
  },
};
