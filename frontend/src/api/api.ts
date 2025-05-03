export const API_URL = "/api";

import { AuthTokens } from "@/types";
import { authAPI } from "./auth";

let isRefreshing = false;
let refreshPromise: Promise<AuthTokens> | null = null;

// Queue of requests to retry after token refresh
const waitingRequests: (() => void)[] = [];

// Process all waiting requests
const processWaitingRequests = () => {
  waitingRequests.forEach((callback) => callback());
  waitingRequests.length = 0;
};

// API client with error handling
export async function apiClient<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  const url = `${API_URL}${endpoint}`;

  // Get token from localStorage as fallback if cookies not available in dev mode
  const token = localStorage.getItem("auth_token");

  // Create headers
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };

  // Add Authorization header if token exists (for dev mode or non-cookie environments)
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  // First attempt
  try {
    const response = await fetch(url, {
      ...options,
      credentials: "include", // For cookies
      headers,
    });

    // If response is successful, return the data
    if (response.ok) {
      return response.json();
    }

    // Handle 401 Unauthorized with token refresh
    if (response.status === 401) {
      // Only try to refresh if not already refreshing
      if (!isRefreshing) {
        isRefreshing = true;
        refreshPromise = authAPI.refreshToken();

        try {
          // Wait for the refresh to complete
          await refreshPromise;

          // Process any waiting requests
          processWaitingRequests();
        } catch (refreshError) {
          // If refresh fails, reject all waiting requests
          waitingRequests.length = 0;
          throw refreshError;
        } finally {
          isRefreshing = false;
          refreshPromise = null;
        }
      } else if (refreshPromise) {
        // Wait for the existing refresh to complete
        try {
          await refreshPromise;
        } catch (error) {
          throw error;
        }
      }

      // Retry the original request with the new token
      const newToken = localStorage.getItem("auth_token");
      const newHeaders = {
        ...headers,
        ...(newToken ? { Authorization: `Bearer ${newToken}` } : {}),
      };

      const retryResponse = await fetch(url, {
        ...options,
        credentials: "include",
        headers: newHeaders,
      });

      if (!retryResponse.ok) {
        const errorData = await retryResponse.json().catch(() => ({}));
        throw new Error(
          errorData.error || `API error: ${retryResponse.status}`,
        );
      }

      return retryResponse.json();
    }

    // For other errors, parse the error message
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || `API error: ${response.status}`);
  } catch (error) {
    console.error(`API request failed for ${url}:`, error);
    throw error;
  }
}
