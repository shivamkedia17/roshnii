export const API_URL = "/api";

export type EndpointParams = {
  endpoint: string;
  options: RequestInit;
  requiresAuth?: boolean;
};

// Simplified API client that relies solely on cookies for authentication
export async function apiClient<T>({
  endpoint = "",
  options = {},
  requiresAuth = undefined,
}: EndpointParams): Promise<T> {
  const url = `${API_URL}${endpoint}`;

  const fetchOptions: RequestInit = {
    ...options,
    headers: {
      ...options.headers, // include all headers specified
      ...(options.body instanceof FormData // don't include content-type for FormData
        ? {}
        : { "Content-Type": "application/json" }),
    },
    credentials: requiresAuth ? "include" : "omit",
  };

  try {
    const response = await fetch(url, fetchOptions);

    // If response is successful, return the data
    if (response.ok) {
      return response.json();
    }

    // For 401 Unauthorized, redirect to login if session expired
    if (response.status === 401) {
      // TODO  - handle 401
      // maybe refresh the token if appropriate
      //
      // or redirect to login endpoint
      //
      // window.dispatchEvent(new CustomEvent("auth:sessionExpired"));
      // You could also redirect to login page directly if appropriate
      // window.location.href = '/api/auth/google/login';
    }

    return (await response.json()) as T;
  } catch (error) {
    console.error(`API request failed for ${url}:`, error);
    throw error;
  }
}

export function ApiHealthCheck() {
  // perform a health-check  using the URL = "/health" endpoint
}
