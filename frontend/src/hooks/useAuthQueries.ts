import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { authAPI } from "@/api/auth";

// Query keys - exported so they can be used elsewhere
export const authKeys = {
  all: ["auth"] as const,
  currentUser: () => [...authKeys.all, "currentUser"] as const,
};

// Get current user hook with improved error handling
export function useCurrentUser() {
  return useQuery({
    queryKey: authKeys.currentUser(),
    queryFn: authAPI.getCurrentUser,
    retry: (failureCount, error) => {
      // Don't retry on 401 Unauthorized - it means we need to login
      if (error instanceof Error && error.message.includes("401")) {
        return false;
      }
      // Retry other errors up to 2 times
      return failureCount < 2;
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}

// Dev login hook
export function useDevLogin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (credentials: { email: string; name: string }) =>
      authAPI.devLogin(credentials),
    onSuccess: () => {
      // Invalidate current user query to refetch
      queryClient.invalidateQueries({ queryKey: authKeys.currentUser() });
    },
  });
}

// Refresh token hook
export function useRefreshToken() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: authAPI.refreshToken,
    onSuccess: () => {
      // After successful token refresh, refetch current user
      queryClient.invalidateQueries({ queryKey: authKeys.currentUser() });
    },
  });
}

// Logout hook
export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      try {
        // Reset auth state immediately to prevent hooks from trying to refresh
        queryClient.setQueryData(authKeys.currentUser(), null);

        // Then perform the actual logout request
        await authAPI.logout();
        return true;
      } catch (error) {
        console.error("Logout API error:", error);
        // Even if the API call fails, still consider the user logged out locally
        return false;
      }
    },
    onSuccess: () => {
      // Reset auth state
      queryClient.setQueryData(authKeys.currentUser(), null);

      // Invalidate all queries to clear cache
      queryClient.invalidateQueries();
    },
  });
}
