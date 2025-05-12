import axios from "axios";
import { API_URL } from "./model";
import { User } from "./model";

export const UserAPI = {
  // axios throws for non-2xx HTTP response codes anyway,
  // so no explicit error handling is needed
  getCurrentUser: async function () {
    const response = await axios.get(`${API_URL}/me`, {
      withCredentials: true,
    });
    console.log(response.status);
    console.log(response.data);
    return response.data as User;
  },
};
