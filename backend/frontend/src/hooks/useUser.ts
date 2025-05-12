// UI Hooks for Current User State
import { User } from "@/api/model";
import { UserAPI } from "@/api/user";
import { useQuery } from "@tanstack/react-query";

// Query keys
export const authKeys = {
  all: ["auth"] as const,
  user: ["auth", "user"] as const,
  loggedIn: ["auth", "loggedIn"] as const,
};

type CurrentUserDetails = {
  isLoading: boolean;
  error?: Error;
  isAuthenticated: boolean;
  currentUser?: User;
};

// Ignore Loading State since we always refetch the query on Mount
export function useGetCurrentUser(): CurrentUserDetails {
  console.log("Getting current user");

  const result = useQuery({
    queryKey: authKeys.user,
    queryFn: UserAPI.getCurrentUser,
    // refetchOnMount: "always",
  });

  // 1. check if loading
  if (result.isLoading) {
    return {
      isLoading: true,
      isAuthenticated: false,
    };
  }

  // 2. on success
  //  a. there must be a user -> return as User, isAuthenticated = true
  if (result.isSuccess) {
    return {
      isLoading: false,
      isAuthenticated: true,
      currentUser: result.data,
    };
  }

  if (result.isError) {
    return {
      isLoading: false,
      isAuthenticated: false,
      error: result.error,
    };
  }

  // For any authorized Requests in axios, redirect User to Login page

  console.error("Please debug your User Fetch function.");
  throw new Error("Debug User query function.");
}
