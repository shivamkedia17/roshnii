import { useEffect, useState } from "react";
import { useAuth } from "@/context/AuthContext";
import "@/css/Auth.css";

export function LogoutPage() {
  const { logout } = useAuth();
  const [status, setStatus] = useState<"loading" | "success" | "error">(
    "loading",
  );
  const [errorMessage, setErrorMessage] = useState("");

  useEffect(() => {
    const performLogout = async () => {
      try {
        await logout();
        setStatus("success");

        // After a successful logout, redirect to login page after a short delay
        setTimeout(() => {
          window.location.reload(); // Force reload to clear any state
        }, 1500);
      } catch (error) {
        console.error("Logout failed:", error);
        setStatus("error");
        setErrorMessage(
          error instanceof Error
            ? error.message
            : "Failed to log out. Please try again.",
        );
      }
    };

    performLogout();
  }, [logout]);

  return (
    <div className="logout-container">
      <div className="logout-card">
        <h2>Logging Out</h2>

        {status === "loading" && (
          <>
            <div className="logout-spinner"></div>
            <p>Logging you out...</p>
          </>
        )}

        {status === "success" && (
          <>
            <div className="logout-success">✓</div>
            <p>You have been successfully logged out.</p>
            <p className="redirect-message">Redirecting to login page...</p>
          </>
        )}

        {status === "error" && (
          <>
            <div className="logout-error">!</div>
            <p className="error-message">{errorMessage}</p>
            <button
              className="retry-button"
              onClick={() => window.location.reload()}
            >
              Try Again
            </button>
          </>
        )}
      </div>
    </div>
  );
}
