import { useAuthContext } from "@/context/AuthContext";
import "@/css/App.css";
import { Loading } from "./common/Loading";
import { useEffect } from "react";

export default function App() {
  const { isAuthenticated, setIsAuthenticated, isLoading, login, logout } =
    useAuthContext();

  useEffect(() => {
    // Add event listener for the authError custom event
    const handleAuthError = () => {
      setIsAuthenticated(false);
    };

    window.addEventListener("authError", handleAuthError);

    // Cleanup listener when component unmounts
    return () => {
      window.removeEventListener("authError", handleAuthError);
    };
  }, [setIsAuthenticated]);

  return (
    <>
      {isLoading ? (
        <Loading />
      ) : isAuthenticated ? (
        <div>
          <h1>Logged In!</h1>
          <button onClick={logout}>Log Out</button>
        </div>
      ) : (
        <button onClick={login}>Login With Google</button>
      )}
    </>
  );
}
