import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";
import {} from "react";

import { authAPI } from "@/services/api";

type AuthContextType = {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: any | null;
  login: () => void;
  logout: () => void;
};

export const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  isLoading: true,
  user: null,
  login: () => {},
  logout: () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [user, setUser] = useState(null);

  useEffect(() => {
    // Check if user is logged in
    async function checkAuth() {
      try {
        setIsLoading(true);

        // Check for token from dev login
        const devToken = localStorage.getItem("auth_token");

        if (devToken) {
          try {
            // Use the token for /me request
            const response = await fetch("/api/me", {
              headers: {
                Authorization: `Bearer ${devToken}`,
              },
            });

            if (response.ok) {
              const userData = await response.json();
              setUser(userData);
              setIsAuthenticated(true);
              return;
            }
          } catch (error) {
            console.error("Dev token validation failed:", error);
            localStorage.removeItem("auth_token");
          }
        }

        // Fall back to regular cookie-based auth check
        const response = await fetch("/api/me", {
          credentials: "include",
        });

        if (response.ok) {
          const userData = await response.json();
          setUser(userData);
          setIsAuthenticated(true);
        }
      } catch (error) {
        console.error("Auth check failed:", error);
      } finally {
        setIsLoading(false);
      }
    }

    checkAuth();
  }, []);

  function login() {
    authAPI.login();
  }

  async function logout() {
    authAPI.logout();
    setIsAuthenticated(false);
    setUser(null);
  }

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        isLoading,
        user,
        login,
        logout,
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
