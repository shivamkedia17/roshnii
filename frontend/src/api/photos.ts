import { apiClient, API_URL } from "./api";

export const photosAPI = {
  // Get all photos
  getPhotos: () => apiClient<any[]>("/images"),

  // Get a single photo
  getPhoto: (id: string) => apiClient<any>(`/image/${id}`),

  // Upload a new photo
  uploadPhoto: async (formData: FormData) => {
    // No longer use localStorage tokens - rely only on cookies
    try {
      const res = await fetch(`${API_URL}/upload`, {
        method: "POST",
        credentials: "include", // This ensures cookies are sent with the request
        // Don't set Content-Type with FormData - browser sets it with boundary
        body: formData,
      });

      if (!res.ok) {
        try {
          const err = await res.json();
          throw new Error(err.error || `Upload failed: ${res.status}`);
        } catch {
          throw new Error(`Upload failed: ${res.status}`);
        }
      }

      return await res.json();
    } catch (error) {
      throw error;
    }
  },

  // Delete a photo
  deletePhoto: (id: string) =>
    apiClient<{ message: string }>(`/image/${id}`, { method: "DELETE" }),
};
