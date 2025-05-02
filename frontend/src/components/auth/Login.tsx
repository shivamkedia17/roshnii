import { useAuth } from "../../context/AuthContext";
import { DevLogin } from "./DevLogin";
import "../../css/Auth.css";

export function Login() {
  // Check if in development mode - you can enhance this with environment checks
  const isDev = true; // or process.env.NODE_ENV === 'development' if configured

  const { login } = useAuth();

  if (isDev) {
    return <DevLogin />;
  }

  // Regular OAuth login for production
  return (
    <div className="login-container">
      <div className="login-card">
        <h1>Roshnii</h1>
        <p>Store and organize your memories</p>
        <button className="google-login-btn" onClick={login}>
          Sign in with Google
        </button>
      </div>
    </div>
  );
}
