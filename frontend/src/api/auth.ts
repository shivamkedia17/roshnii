import { apiClient, API_URL } from "./api";
import { UserInfo } from "@/types";
import { AuthTokens } from "@/types";

export const authAPI = {
  login: () => (window.location.href = "/api/auth/google/login"),

  // Enhanced logout with return value
  logout: async (): Promise<void> => {
    try {
      await apiClient("/auth/google/logout", { method: "POST" });
      // Clear any local storage items
      localStorage.removeItem("auth_token");
      localStorage.removeItem("refresh_token");
    } catch (error) {
      console.error("Logout error:", error);
      // Still clear local storage even if API call fails
      localStorage.removeItem("auth_token");
      localStorage.removeItem("refresh_token");
      throw error;
    }
  },

  // Get current user with better error handling
  getCurrentUser: async (): Promise<UserInfo> => {
    try {
      return await apiClient<UserInfo>("/me");
    } catch (error) {
      console.error("Error fetching current user:", error);
      throw error;
    }
  },

  // New refresh token function
  refreshToken: async (): Promise<AuthTokens> => {
    try {
      // Use the refresh_token cookie or token from localStorage
      const refreshToken = localStorage.getItem("refresh_token");
      const headers: HeadersInit = {};

      if (refreshToken) {
        headers["Authorization"] = `Bearer ${refreshToken}`;
      }

      const response = await fetch(`${API_URL}/auth/google/refresh`, {
        method: "POST",
        credentials: "include", // For cookies
        headers,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(
          errorData.error || `Refresh failed: ${response.status}`,
        );
      }

      const data = await response.json();

      // In dev mode, the API returns tokens in the response
      // Store them in localStorage as fallback
      if (data.token) {
        localStorage.setItem("auth_token", data.token);
      }
      if (data.refresh_token) {
        localStorage.setItem("refresh_token", data.refresh_token);
      }

      return {
        token: data.token || "",
        refreshToken: data.refresh_token,
        expiresIn: data.expires_in,
      };
    } catch (error) {
      console.error("Token refresh error:", error);
      throw error;
    }
  },

  // Dev login with token storage
  devLogin: async (credentials: {
    email: string;
    name: string;
  }): Promise<AuthTokens> => {
    const response = await fetch(`${API_URL}/auth/dev/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(credentials),
      credentials: "include", // For cookies
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `Login failed: ${response.status}`);
    }

    const data = await response.json();

    // Store tokens in localStorage for dev mode
    if (data.token) {
      localStorage.setItem("auth_token", data.token);
    }
    if (data.refresh_token) {
      localStorage.setItem("refresh_token", data.refresh_token);
    }

    return {
      token: data.token || "",
      refreshToken: data.refresh_token,
      expiresIn: data.expires_in,
    };
  },
};
