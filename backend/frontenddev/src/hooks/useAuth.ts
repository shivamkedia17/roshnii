// src/hooks/useAuth.ts
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AuthAPI } from "@/api/auth";
import { UserAPI } from "@/api/user";

// Query keys
export const authKeys = {
  all: ["auth"],
  user: ["auth", "user"],
  loggedIn: ["auth", "loggedIn"],
};

// Mutation hook for Google OAuth login
export function useLogin() {
  return useMutation({
    mutationFn: AuthAPI.login,
    // onSuccess: useCurrentUser,
    // onError: (error) => {
    //   console.error("Login error", error);
    // },
    // throwOnError: true,
  });
}

// Mutation hook for logging out
export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: AuthAPI.logout,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: authKeys.all });
    },
    onError: (error) => {
      console.error("Logout error", error);
    },
    throwOnError: true,
  });
}

export function useCurrentUser() {
  console.log("Calling get current user");
  return useQuery({
    queryKey: authKeys.user,
    queryFn: UserAPI.getCurrentUser,
    refetchOnMount: "always",
  });
}
