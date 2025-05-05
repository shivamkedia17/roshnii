import { API_URL, apiClient, EndpointParams } from "./api";
import { ServerMessage } from "./model";

export const AuthAPI = {
  baseEndpoint: "/auth/google",

  login: async function () {
    try {
      let url = `${API_URL}/auth/google/login`;
      let fetchOptions: RequestInit = {
        method: "GET",
        redirect: "follow",
      };

      let response = await fetch(url, fetchOptions);

      if (!response.ok) {
        throw new Error(`Login Response not OK:\n${await response.json()}`);
      }

      const { auth_url } = (await response.json()) as { auth_url: string };

      window.location.href = auth_url;

      // let finalURL = new URL(response.url);
      // const error = finalURL.searchParams.get("error");
      // const code = finalURL.searchParams.get("code");
      // const state = finalURL.searchParams.get("state");

      // if (!code || !state) {
      //   throw new Error("Missing required OAuth parameters");
      // }

      // console.log("OAuth Params: ", error, code, state);

      // url = `${API_URL}/auth/google/callback`;
      // fetchOptions = {
      //   method: "GET",
      //   redirect: "follow",
      //   credentials: "include",
      // };

      // response = await fetch(url, fetchOptions);

      // if (!response.ok) {
      //   throw new Error(`Callback Response not OK:\n${await response.json()}`);
      // }

      // console.log(await response.json());
    } catch (err) {
      console.error("Error logging in: ", err);
      throw err;
    }
  },

  // Logs out the current user
  logout: async function () {
    const params: EndpointParams = {
      endpoint: `/auth/google/logout`,
      includeCookies: true,
      options: {
        method: "POST",
      },
    };

    return await apiClient<{ message: string }>(params);
  },

  // Refreshes the authentication token
  refreshToken: async function () {
    try {
      const url = `${API_URL}/auth/google/refresh`;
      const fetchOptions: RequestInit = {
        method: "GET",
        credentials: "include", // Important: includes cookies in the request
        headers: {
          Accept: "application/json",
        },
      };

      const response = await fetch(url, fetchOptions);

      if (!response.ok) {
        throw new Error(
          `Refresh token response not OK: ${await response.text()}`,
        );
      }

      // The server sets the new access token as a cookie automatically
      // We just need to return the response data
      return (await response.json()) as ServerMessage;
    } catch (err) {
      console.error("Error refreshing token: ", err);
      throw err;
    }
  },
};
