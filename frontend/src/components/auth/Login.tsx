import { useState } from "react";
import { useAuth } from "@/context/AuthContext";
import { DevLogin } from "./DevLogin";
import "@/css/Auth.css";

export function Login() {
  const [loginError, setLoginError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // Check if in development mode
  const isDev = import.meta.env.NODE_ENV === "development" || true; // For testing, default to true

  const { login } = useAuth();

  const handleOAuthLogin = () => {
    setIsLoading(true);
    setLoginError(null);

    try {
      login();
    } catch (error) {
      setLoginError("Failed to redirect to login page. Please try again.");
      setIsLoading(false);
    }
  };

  if (isDev) {
    return <DevLogin />;
  }

  // Regular OAuth login for production
  return (
    <div className="login-container">
      <div className="login-card">
        <h1>Roshnii</h1>
        <p>Store and organize your memories</p>

        {loginError && <div className="login-error">{loginError}</div>}

        <button
          className={`google-login-btn ${isLoading ? "loading" : ""}`}
          onClick={handleOAuthLogin}
          disabled={isLoading}
        >
          {isLoading ? "Redirecting..." : "Sign in with Google"}
        </button>
      </div>
    </div>
  );
}
