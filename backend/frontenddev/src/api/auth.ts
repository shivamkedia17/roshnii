import axios from "axios";
import axiosInstance from "./api";
import { API_URL } from "./model";

export const AuthAPI = {
  login: async function () {
    try {
      const response = await axios.get(`${API_URL}/auth/google/login`);

      if (response.data && response.data.auth_url) {
        window.location.href = response.data.auth_url;
      } else {
        throw new Error("Invalid authentication URL received");
      }
    } catch (err) {
      console.error("Error logging in: ", err);
      throw err;
    }
  },

  logout: async function () {
    const response = await axiosInstance.post("/auth/google/logout");
    return response.data;
  },
};
