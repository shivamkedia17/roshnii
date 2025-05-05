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
// FIXME this clearly doesn't work properly
export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: authAPI.logout,
    onSuccess: () => {
      // Reset auth state and clear cache on successful logout
      queryClient.setQueryData(authKeys.currentUser(), null);
      queryClient.clear(); // Clear the entire cache
    },
  });
}
