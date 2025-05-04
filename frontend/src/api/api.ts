export const API_URL = "/api";

// Simplified API client that relies solely on cookies for authentication
export async function apiClient<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  const url = `${API_URL}${endpoint}`;

  // Create headers with default content type
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };

  try {
    const response = await fetch(url, {
      ...options,
      credentials: "include", // This ensures cookies are sent with every request
      headers,
    });

    // If response is successful, return the data
    if (response.ok) {
      return response.json();
    }

    // For 401 Unauthorized, redirect to login if session expired
    if (response.status === 401) {
      // Optional: Dispatch an event that the app can listen for to show login prompt
      window.dispatchEvent(new CustomEvent("auth:sessionExpired"));
      // You could also redirect to login page directly if appropriate
      // window.location.href = '/api/auth/google/login';
    }

    // For other errors, parse the error message
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || `API error: ${response.status}`);
  } catch (error) {
    console.error(`API request failed for ${url}:`, error);
    throw error;
  }
}
