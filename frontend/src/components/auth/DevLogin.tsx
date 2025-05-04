import { useState, useEffect } from "react";
import { useDevLogin } from "@/hooks/useAuthQueries";

export function DevLogin() {
  // Add state to detect if we're in development mode
  const [email, setEmail] = useState("testuser@example.com");
  const [name, setName] = useState("Test User");
  const [loggingIn, setLoggingIn] = useState(false);
  const [isDev, setIsDev] = useState(false); // New state

  const devLoginMutation = useDevLogin();

  // Add effect to check for development environment
  useEffect(() => {
    const isDevEnvironment =
      import.meta.env.DEV ||
      window.location.hostname === "localhost" ||
      window.location.hostname === "127.0.0.1";

    setIsDev(isDevEnvironment);
  }, []);

  const handleDevLogin = async (e: React.FormEvent) => {
    e.preventDefault();

    // Add safety check
    if (!isDev) {
      console.error("Dev login is only available in development mode");
      return;
    }

    setLoggingIn(true);

    try {
      // Simplified - no more localStorage token handling here
      await devLoginMutation.mutateAsync({ email, name });
    } catch (error) {
      console.error("Dev login error:", error);
    } finally {
      setLoggingIn(false);
    }
  };

  if (!isDev) {
    return (
      <div className="dev-login-container">
        <div className="login-card">
          <h1>Development Mode Only</h1>
          <p>
            This login method is only available in development environments.
          </p>
          <button
            className="google-login-btn"
            onClick={() => (window.location.href = "/api/auth/google/login")}
          >
            Sign in with Google
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="dev-login-container">
      <div className="login-card">
        {/* Add prominent development warning */}
        <div className="dev-warning">DEVELOPMENT MODE ONLY</div>
        <h1>Roshnii</h1>
        <p>Development Login</p>

        {/* Error handling remains the same */}
        {devLoginMutation.error && (
          <div className="error">
            {devLoginMutation.error instanceof Error
              ? devLoginMutation.error.message
              : "Login failed. Please try again."}
          </div>
        )}

        <form onSubmit={handleDevLogin}>
          <div className="form-group">
            <label htmlFor="email">Email:</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={loggingIn || devLoginMutation.isPending}
            />
          </div>

          <div className="form-group">
            <label htmlFor="name">Name:</label>
            <input
              type="text"
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              disabled={loggingIn || devLoginMutation.isPending}
            />
          </div>

          <button
            type="submit"
            disabled={loggingIn || devLoginMutation.isPending || !email}
            className={loggingIn || devLoginMutation.isPending ? "loading" : ""}
          >
            {loggingIn || devLoginMutation.isPending
              ? "Logging in..."
              : "Dev Login"}
          </button>
        </form>

        <div className="oauth-option">
          <hr />
          <p>Or use OAuth (if configured):</p>
          <button
            className="google-login-btn"
            onClick={() => (window.location.href = "/api/auth/google/login")}
            disabled={loggingIn || devLoginMutation.isPending}
          >
            Sign in with Google
          </button>
        </div>
      </div>
    </div>
  );
}
