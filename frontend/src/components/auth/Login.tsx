import "../../css/Auth.css";

export function Login() {
  const handleGoogleLogin = () => {
    // Implement OAuth login with Google
    window.location.href = "/api/login"; // Redirect to backend OAuth endpoint
  };

  return (
    <div className="login-container">
      <div className="login-card">
        <h1>Roshnii</h1>
        <p>Store and organize your memories</p>
        <button className="google-login-btn" onClick={handleGoogleLogin}>
          Sign in with Google
        </button>
      </div>
    </div>
  );
}
