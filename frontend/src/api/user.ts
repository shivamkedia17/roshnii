import { apiClient } from "./api";
import { User } from "./model";

export const UserAPI = {
  baseEndpoint: "/me",

  getCurrentUser: function () {
    const params = {
      endpoint: this.baseEndpoint,
      requiresAuth: true,
      options: {
        method: "GET",
      },
    };

    return apiClient<User>(params);
  },
};
