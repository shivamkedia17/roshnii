import { apiClient, API_URL } from "./api";

export const photosAPI = {
  // Get all photos
  getPhotos: () => apiClient<any[]>("/images"),

  // Get a single photo
  getPhoto: (id: string) => apiClient<any>(`/image/${id}`),

  // Upload a new photo
  uploadPhoto: (formData: FormData) => {
    const token = localStorage.getItem("auth_token");
    const headers: HeadersInit = {};

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    return fetch(`${API_URL}/upload`, {
      method: "POST",
      credentials: "include",
      headers,
      body: formData,
    }).then((res) => {
      if (!res.ok) {
        return res
          .json()
          .then((err) => {
            throw new Error(err.error || `Upload failed: ${res.status}`);
          })
          .catch(() => {
            throw new Error(`Upload failed: ${res.status}`);
          });
      }
      return res.json();
    });
  },

  // Delete a photo
  deletePhoto: (id: string) =>
    apiClient<{ message: string }>(`/image/${id}`, { method: "DELETE" }),
};
