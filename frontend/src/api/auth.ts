import { apiClient, API_URL } from "./api";
import { UserInfo } from "@/types";

// Environment detection
const isDevelopment =
  import.meta.env.DEV ||
  window.location.hostname === "localhost" ||
  window.location.hostname === "127.0.0.1";

export const authAPI = {
  // Redirect to OAuth login
  login: () => {
    window.location.href = "/api/auth/google/login";
  },

  // Simple logout that relies on backend to clear cookies
  logout: async (): Promise<void> => {
    try {
      await fetch("/api/auth/google/logout", {
        method: "POST",
        credentials: "include",
        headers: {
          "Cache-Control": "no-cache, no-store",
          Pragma: "no-cache",
        },
      });
    } catch (error) {
      console.error("Logout error:", error);
      // Don't throw here, just log the error
      // This makes the logout more fault-tolerant
    }
  },

  // Get current user
  getCurrentUser: async (): Promise<UserInfo> => {
    try {
      return await apiClient<UserInfo>("/me");
    } catch (error) {
      console.error("Error fetching current user:", error);
      throw error;
    }
  },

  // Refresh token - just calls the endpoint, backend handles cookie updates
  refreshToken: async (): Promise<{ message: string }> => {
    try {
      return await apiClient<{ message: string }>("/auth/google/refresh", {
        method: "POST",
        credentials: "include",
      });
    } catch (error) {
      console.error("Token refresh error:", error);
      throw error;
    }
  },

  // Development-only login feature
  // This should ONLY be available in development mode
  devLogin: async (credentials: {
    email: string;
    name: string;
  }): Promise<void> => {
    if (!isDevelopment) {
      console.error("Dev login is only available in development mode");
      throw new Error("Dev login not available");
    }

    const response = await fetch(`${API_URL}/auth/dev/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(credentials),
      credentials: "include", // Still use cookies even for dev login
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `Login failed: ${response.status}`);
    }

    return response.json();
  },
};
