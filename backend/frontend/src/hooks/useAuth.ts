// Hooks to Set and Unset Cookies that deal with UserAuth
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AuthAPI } from "@/api/auth";
import { authKeys } from "./useUser";

// Mutation hook to call OAuth login endpoint
export function useLogin() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    // The webapp is context switched out due to window navigation,
    // hence only cookies are set.
    mutationFn: AuthAPI.login,
    onSuccess: () => {
      // Invalidate the queries that store current user data, so that they are refetched.
      queryClient.invalidateQueries({ queryKey: authKeys.all });
    },
    onError: (error) => {
      console.error("Login error", error);
      // TODO? could not log user in
    },
    throwOnError: true,
  });

  return mutation.mutate;
}

// Mutation hook for logging out
export function useLogout() {
  const queryClient = useQueryClient();

  // backend server returns a message only in the absence of any errors
  // backend server clear cookies on successful request
  const mutation = useMutation({
    mutationFn: AuthAPI.logout,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: authKeys.all });
    },
    onError: (error) => {
      console.error("Logout error", error);
      // TODO? could not log user out
    },
    throwOnError: true,
  });

  return mutation.mutate;
}
