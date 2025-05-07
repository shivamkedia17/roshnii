import axios from "axios";
import { API_URL } from "./model";
import { User } from "./model";

export const UserAPI = {
  getCurrentUser: async function () {
    const response = await axios.get(`${API_URL}/me`, {
      withCredentials: true,
    });
    console.log(response.status);
    console.log(response.data);
    // const response = await axiosInstance.get("/me");
    return response.data as User;
  },
};
