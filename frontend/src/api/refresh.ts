import axios from "axios";
import { API_URL } from "./model";
import { ServerMessage } from "./model";

export const RefreshAuthAPI = {
  async refreshToken() {
    try {
      const response = await axios.get(`${API_URL}/auth/google/refresh`, {
        withCredentials: true,
      });
      return response.data as ServerMessage;
    } catch (err) {
      console.error("Error refreshing token: ", err);
      throw err;
    }
  },
};
