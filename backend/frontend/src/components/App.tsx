import "@/css/App.css";
import { useGetCurrentUser } from "@/hooks/useUser";
import { useEffect } from "react";
import { Login } from "./auth/Login";
import { Loading } from "./common/Loading";
import { MainLayout } from "./layout/MainLayout";

export default function App() {
  // TODO substitute useEffect for a custom hook that listens for authError and refreshError
  // 1. redirect to (Login) for RefreshError event
  // 2. Show Error Boundary for AuthError event

  useEffect(() => {
    // Add event listener for the authError custom event
    const handleAuthError = () => {
      console.log("Some 401 Auth Error Occurred.");
    };

    const handleRefreshError = () => {
      console.log("Some Refresh Error occurred.");
    };

    window.addEventListener("authError", handleAuthError);
    window.addEventListener("refreshError", handleRefreshError);

    // Cleanup listener when component unmounts
    return () => {
      window.removeEventListener("authError", handleAuthError);
      window.removeEventListener("refreshError", handleRefreshError);
    };
  }, []);

  const userResult = useGetCurrentUser();

  if (userResult.isLoading) {
    return <Loading />;
  } else {
    if (!userResult.isAuthenticated) {
      return <Login />;
    } else {
      return <MainLayout />;
    }
  }
}
