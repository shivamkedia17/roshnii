import {
  createContext,
  useContext,
  ReactNode,
  useCallback,
  useEffect,
} from "react";
import {
  useCurrentUser,
  useRefreshToken,
  useLogout,
  authKeys,
} from "@/hooks/useAuthQueries";
import { AuthContextType, UserInfo } from "@/types";
import { useQueryClient } from "@tanstack/react-query";

export const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  isLoading: true,
  user: null,
  login: () => {},
  logout: async () => true,
  refreshToken: async () => false,
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const { data: user, isLoading, error: userError, isError } = useCurrentUser();
  const refreshMutation = useRefreshToken();
  const logoutMutation = useLogout();
  const queryClient = useQueryClient();

  // Function to refresh token
  const refreshToken = useCallback(async (): Promise<boolean> => {
    try {
      // Don't attempt refresh if we've explicitly logged out
      if (
        isError &&
        userError instanceof Error &&
        userError.message.includes("logged out")
      ) {
        return false;
      }

      await refreshMutation.mutateAsync();
      return true;
    } catch (error) {
      console.error("Token refresh error:", error);
      return false;
    }
  }, [refreshMutation, isError, userError]);

  // Listen for session expired events
  useEffect(() => {
    const handleSessionExpired = () => {
      refreshToken().catch((err) => {
        console.error("Auto-refresh failed:", err);
      });
    };

    window.addEventListener("auth:sessionExpired", handleSessionExpired);
    return () => {
      window.removeEventListener("auth:sessionExpired", handleSessionExpired);
    };
  }, [refreshToken]);

  // Setup automatic refresh
  useEffect(() => {
    // If auth error is due to expired token, try to refresh
    if (
      isError &&
      userError instanceof Error &&
      userError.message.includes("expired")
    ) {
      refreshToken().catch((err) => {
        console.error("Auto-refresh failed:", err);
      });
    }
  }, [isError, userError, refreshToken]);

  // Redirect to Google login
  const login = useCallback(() => {
    window.location.href = "/api/auth/google/login";
  }, []);

  // Function to handle logout
  const logout = useCallback(async () => {
    try {
      // First set authentication state to logged out to prevent refresh attempts
      // This prevents the infinite loop
      queryClient.setQueryData(authKeys.currentUser(), null);

      // Then perform the actual logout API call
      return await logoutMutation.mutateAsync();
    } catch (error) {
      console.error("Logout error:", error);
      return false;
    }
  }, [queryClient, logoutMutation]);

  // Determine authentication status
  const isAuthenticated = !!user && !isError;

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        isLoading,
        user: user as UserInfo | null,
        login,
        logout,
        refreshToken,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used within an AuthContextProvider");
  }

  return context;
};
