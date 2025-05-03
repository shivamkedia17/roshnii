import { useState } from "react";
import { useDevLogin } from "@/hooks/useAuthQueries";

export function DevLogin() {
  const [email, setEmail] = useState("testuser@example.com");
  const [name, setName] = useState("Test User");
  const [loggingIn, setLoggingIn] = useState(false);

  const devLoginMutation = useDevLogin();

  const handleDevLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoggingIn(true);

    try {
      await devLoginMutation.mutateAsync({ email, name });
    } catch (error) {
      console.error("Dev login error:", error);
    } finally {
      setLoggingIn(false);
    }
  };

  return (
    <div className="dev-login-container">
      <div className="login-card">
        <h1>Roshnii</h1>
        <p>Development Login</p>

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
