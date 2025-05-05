import { apiClient, EndpointParams, API_URL } from "./api";

export const AuthAPI = {
  baseEndpoint: "/auth/google",

  // Returns the Google OAuth login URL (redirect user to this URL)
  getLoginURL: function () {
    return `${API_URL}${this.baseEndpoint}/login`;
  },

  // Not directly called - backend redirects to this URL after Google auth
  getCallbackURL: function () {
    return `${API_URL}${this.baseEndpoint}/callback`;
  },

  // Refreshes the authentication token
  refreshToken: function () {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/refresh`,
      requiresAuth: true,
      options: {
        method: "POST",
      },
    };

    return apiClient<{ message: string }>(params);
  },

  // Logs out the current user
  logout: function () {
    const params: EndpointParams = {
      endpoint: `${this.baseEndpoint}/logout`,
      requiresAuth: true,
      options: {
        method: "POST",
      },
    };

    return apiClient<{ message: string }>(params);
  },
};
