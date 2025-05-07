import { createContext, useContext, useState, ReactNode } from "react";

type MockAuthContextType = {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: any;
  login: () => void;
  logout: () => void;
};

const MockAuthContext = createContext<MockAuthContextType>({
  isAuthenticated: true, // Always authenticated in mock mode
  isLoading: false,
  user: { id: 1, email: "mock@example.com", name: "Mock User" },
  login: () => {},
  logout: () => {},
});

export function MockAuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(true);
  const [user] = useState({
    id: 1,
    email: "mock@example.com",
    name: "Mock User",
  });

  const login = () => setIsAuthenticated(true);
  const logout = () => setIsAuthenticated(false);

  return (
    <MockAuthContext.Provider
      value={{
        isAuthenticated,
        isLoading: false,
        user,
        login,
        logout,
      }}
    >
      {children}
    </MockAuthContext.Provider>
  );
}

export const useMockAuth = () => {
  return useContext(MockAuthContext);
};
