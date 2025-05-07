import axios, { AxiosInstance } from "axios";
import { API_URL } from "./model";
import { RefreshAuthAPI } from "./refresh";

// Create a configured axios instance for the application
const axiosInstance: AxiosInstance = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true, // Include cookies by default
});

// Intercept 401 unauthorized errors to handle token refresh
axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Check if the error is a 401 unauthorized error
    if (error.response?.status === 401) {
      // Check if the error message indicates token refresh is needed
      if (
        error.response.data?.message?.includes("expired token") &&
        !originalRequest._retry
      ) {
        originalRequest._retry = true;

        try {
          // Attempt to refresh the token
          await RefreshAuthAPI.refreshToken();

          // Retry the original request with the new token
          return axiosInstance(originalRequest);
        } catch (refreshError) {
          // If refresh fails, dispatch an auth error event
          window.dispatchEvent(
            new CustomEvent("authError", {
              detail: { message: "Authentication failed" },
            }),
          );
          return Promise.reject(refreshError);
        }
      } else {
        // For other 401 errors, dispatch auth error event
        window.dispatchEvent(
          new CustomEvent("authError", {
            detail: { message: "Authentication failed" },
          }),
        );
      }
    }

    return Promise.reject(error);
  },
);

export function ApiHealthCheck() {
  return axios.get("/health");
}

// Export the axios instance for use in other API modules
export default axiosInstance;
