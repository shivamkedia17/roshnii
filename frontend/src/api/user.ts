import { apiClient, EndpointParams } from "./api";
import { User } from "./model";

export const UserAPI = {
  baseEndpoint: `/me`,

  getCurrentUser: async function () {
    const params: EndpointParams = {
      endpoint: this.baseEndpoint,
      includeCookies: true,
      options: {
        method: "GET",
      },
    };

    const user: User = await apiClient<User>(params);
    return user;
  },
};
