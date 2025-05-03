import {
  createContext,
  useContext,
  ReactNode,
  useCallback,
  useEffect,
} from "react";
import {
  useCurrentUser,
  useLogout,
  useRefreshToken,
} from "@/hooks/useAuthQueries";
import { AuthContextType, UserInfo } from "@/types";

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

  const logoutMutation = useLogout();
  const refreshMutation = useRefreshToken();

  // Function to refresh token
  const refreshToken = useCallback(async (): Promise<boolean> => {
    try {
      await refreshMutation.mutateAsync();
      return true;
    } catch (error) {
      console.error("Token refresh error:", error);
      return false;
    }
  }, [refreshMutation]);

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
      await logoutMutation.mutateAsync();
      return true;
    } catch (error) {
      console.error("Logout error:", error);
      return false;
    }
  }, [logoutMutation]);

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
