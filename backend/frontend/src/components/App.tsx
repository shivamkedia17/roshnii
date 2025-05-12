import "@/css/App.css";
import { Loading } from "./common/Loading";
// import { useEffect } from "react";
import { useGetCurrentUser } from "@/hooks/useUser";
import { TempLogout } from "./TempLogout";
import { TempLogin } from "./TempLogin";

export default function App() {
  // TODO substitute useEffect for a custom hook that listens for authError and refreshError
  // 1. redirect to (Login) for RefreshError event
  // 2. Show Error Boundary for AuthError event

  // useEffect(() => {
  //   // Add event listener for the authError custom event
  //   const handleAuthError = () => {
  //     // TODO
  //   };

  //   const handleRefreshError = () => {
  //     // TODO
  //   };

  //   window.addEventListener("authError", handleAuthError);
  //   window.addEventListener("refreshError", handleRefreshError);

  //   // Cleanup listener when component unmounts
  //   return () => {
  //     window.removeEventListener("authError", handleAuthError);
  //     window.removeEventListener("refreshError", handleRefreshError);
  //   };
  // }, []);

  const userResult = useGetCurrentUser();

  if (userResult.isLoading) {
    return <Loading />;
  } else {
    if (userResult.isAuthenticated) {
      return <TempLogout />;
    } else {
      return <TempLogin />;
    }
  }
}
