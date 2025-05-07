import { useState, useEffect } from "react";
import { useAuth } from "@/context/AuthContext";
import "@/css/SessionExpired.css";

type SessionExpiredModalProps = {
  onClose: () => void;
};

export function SessionExpiredModal({ onClose }: SessionExpiredModalProps) {
  const { refreshToken, login } = useAuth();
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [refreshError, setRefreshError] = useState<string | null>(null);

  // Auto-refresh on mount
  useEffect(() => {
    handleRefresh();
  }, []);

  const handleRefresh = async () => {
    setIsRefreshing(true);
    setRefreshError(null);

    try {
      const success = await refreshToken();
      if (success) {
        onClose(); // Close modal on success
      } else {
        setRefreshError(
          "Could not refresh your session. Please try again or log in.",
        );
      }
    } catch (error) {
      console.error("Session refresh error:", error);
      setRefreshError(
        error instanceof Error
          ? error.message
          : "Failed to refresh your session.",
      );
    } finally {
      setIsRefreshing(false);
    }
  };

  const handleLogin = () => {
    login(); // Redirect to login page
  };

  return (
    <div className="session-expired-overlay">
      <div className="session-expired-modal">
        <h2>Session Expired</h2>
        <p>Your session has expired due to inactivity.</p>

        {refreshError && <div className="refresh-error">{refreshError}</div>}

        <div className="session-buttons">
          <button
            className="refresh-button"
            onClick={handleRefresh}
            disabled={isRefreshing}
          >
            {isRefreshing ? "Refreshing..." : "Refresh Session"}
          </button>

          <button
            className="login-button"
            onClick={handleLogin}
            disabled={isRefreshing}
          >
            Log In Again
          </button>
        </div>
      </div>
    </div>
  );
}
