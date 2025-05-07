import { User } from "@/api/model";
import { useCurrentUser, useLogin, useLogout } from "@/hooks/useAuth";
import { createContext, ReactNode, useContext, useState } from "react";

type AuthContextProps = {
  children: ReactNode;
};

type AuthContextType = {
  user?: User;
  isAuthenticated: boolean;
  setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>;
  isLoading: boolean;
  error: Error | null;
  login: () => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: AuthContextProps) {
  // Use the TanStack Query hooks

  const loginMutation = useLogin();
  const logoutMutation = useLogout();

  // Helper functions to expose mutations more cleanly
  const login = async () => {
    await loginMutation.mutateAsync();
  };

  const logout = async () => {
    await logoutMutation.mutateAsync();
  };

  const { data: user, isLoading, error } = useCurrentUser();
  const [isAuthenticated, setIsAuthenticated] = useState(!!user);
  console.log(user);
  console.log(isAuthenticated);
  // Determine if the user is authenticated

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated,
        setIsAuthenticated,
        error,
        isLoading,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuthContext() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error(
      "useAuthContext must be used within an AuthProvider component.",
    );
  }

  return context;
}
