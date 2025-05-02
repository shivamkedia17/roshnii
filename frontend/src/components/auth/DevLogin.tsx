// src/components/auth/DevLogin.tsx
import { useState } from "react";
import { useAuth } from "../../context/AuthContext";

export function DevLogin() {
  const [email, setEmail] = useState("testuser@example.com");
  const [name, setName] = useState("Test User");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const { login } = useAuth();

  const handleDevLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await fetch("/api/auth/dev/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, name }),
      });

      if (!response.ok) {
        throw new Error("Login failed");
      }

      const data = await response.json();

      // Store the token in localStorage for use with API requests
      localStorage.setItem("auth_token", data.token);

      // Refresh the page or update auth state
      window.location.reload();
    } catch (err) {
      setError("Login failed. Please try again.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="dev-login-container">
      <div className="login-card">
        <h1>Roshnii</h1>
        <p>Development Login</p>

        {error && <div className="error">{error}</div>}

        <form onSubmit={handleDevLogin}>
          <div className="form-group">
            <label htmlFor="email">Email:</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
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
            />
          </div>

          <button type="submit" disabled={loading}>
            {loading ? "Logging in..." : "Dev Login"}
          </button>
        </form>

        <div className="oauth-option">
          <hr />
          <p>Or use OAuth (if configured):</p>
          <button
            className="google-login-btn"
            onClick={login}
            disabled={loading}
          >
            Sign in with Google
          </button>
        </div>
      </div>
    </div>
  );
}
