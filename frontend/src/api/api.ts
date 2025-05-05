import { useAuthContext } from "@/context/AuthContext";
import { AuthAPI } from "./auth";
import { ServerMessage } from "./model";

export const API_URL = "/api";

export type EndpointParams = {
  endpoint: string;
  options: RequestInit;
  includeCookies?: boolean;
};

// Simplified API client that relies solely on cookies for authentication
export async function apiClient<T>({
  endpoint = "",
  options = {},
  includeCookies = undefined,
}: EndpointParams): Promise<T | undefined> {
  const url = `${API_URL}${endpoint}`;

  const fetchOptions: RequestInit = {
    ...options,
    headers: {
      ...options.headers, // include all headers specified
      ...(options.body instanceof FormData // don't include content-type for FormData
        ? {}
        : { "Content-Type": "application/json" }),
    },
    credentials: includeCookies ? "include" : "omit",
  };

  try {
    const response = await fetch(url, fetchOptions);

    // If response is successful, return the data
    if (response.ok) {
      console.log(response.json());
      return (await response.json()) as T;
    }

    // For 401 Unauthorized, redirect to login if unauthorized after trying to refresh JWT
    if (response.status === 401) {
      const body = (await response.json()) as ServerMessage;

      // messages copied from backend
      if (body.message.includes("please refresh your token")) {
        // try refreshing the token
        const attempt = await AuthAPI.refreshToken();
        if (!attempt || attempt.message != "Token refreshed successfully") {
          const { setIsAuthenticated } = useAuthContext();
          setIsAuthenticated(false);
        }
      }

      // or redirect to login endpoint
      const { setIsAuthenticated } = useAuthContext();
      setIsAuthenticated(false);
    }
  } catch (error) {
    console.error(`API request failed for ${url}:`, error);
    throw error;
  }
}

export function ApiHealthCheck() {
  // perform a health-check  using the URL = "/health" endpoint
}
