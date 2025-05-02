const API_URL = "/api";

// API client with error handling
async function apiClient<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  // Get token from localStorage if available (dev mode)
  const token = localStorage.getItem("auth_token");

  // Create a new headers object
  const headerObj: HeadersInit = {
    "Content-Type": "application/json",
  };

  // Add Authorization header if token exists
  if (token) {
    headerObj["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    credentials: "include", // For cookies
    headers: {
      ...headerObj,
      ...(options.headers || {}),
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.error || `API error: ${response.status}`);
  }

  return response.json();
}

export const authAPI = {
  login: () => (window.location.href = "/api/auth/google/login"),
  logout: () => apiClient("/auth/google/logout", { method: "POST" }),
  getCurrentUser: () => apiClient<any>("/me"),
};

export const photosAPI = {
  getPhotos: () => apiClient<any[]>("/images"),
  getPhoto: (id: string) => apiClient<any>(`/image/${id}`),
  uploadPhoto: (formData: FormData) =>
    fetch(`${API_URL}/upload`, {
      method: "POST",
      credentials: "include",
      body: formData,
    }).then((res) => {
      if (!res.ok) throw new Error("Upload failed");
      return res.json();
    }),
};
