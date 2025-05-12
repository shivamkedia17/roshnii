import { useLogin } from "@/hooks/useAuth";
import "@/css/Auth.css";

export function Login() {
  const login = useLogin();

  return (
    <div className="login-container">
      <div className="login-card">
        <h1>Roshnii</h1>
        <p>Store and organize your memories</p>
        <button className="google-login-btn" onClick={() => login()}>
          Login With Google
        </button>
      </div>
    </div>
  );
}
